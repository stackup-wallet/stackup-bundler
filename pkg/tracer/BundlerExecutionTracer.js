var tracer = {
  reverts: [],
  validationOOG: false,
  executionOOG: false,
  executionGasLimit: 0,

  _depth: 0,
  _executionGasStack: [],
  _defaultGasItem: { used: 0, required: 0 },
  _marker: 0,
  _validationMarker: 1,
  _executionMarker: 3,
  _userOperationEventTopics0:
    "0x49628fd1471006c1482da88028e9ce4dbb080b815c9b0344d39e5a8e6ec1419f",

  _isValidation: function () {
    return (
      this._marker >= this._validationMarker &&
      this._marker < this._executionMarker
    );
  },

  _isExecution: function () {
    return this._marker === this._executionMarker;
  },

  _isUserOperationEvent: function (log) {
    var topics0 = "0x" + log.stack.peek(2).toString(16);
    return topics0 === this._userOperationEventTopics0;
  },

  _setUserOperationEvent: function (opcode, log) {
    var count = parseInt(opcode.substring(3));
    var ofs = parseInt(log.stack.peek(0).toString());
    var len = parseInt(log.stack.peek(1).toString());
    var topics = [];
    for (var i = 0; i < count; i++) {
      topics.push("0x" + log.stack.peek(2 + i).toString(16));
    }
    var data = toHex(log.memory.slice(ofs, ofs + len));
    this.userOperationEvent = {
      topics: topics,
      data: data,
    };
  },

  fault: function fault(log, db) {},
  result: function result(ctx, db) {
    return {
      reverts: this.reverts,
      validationOOG: this.validationOOG,
      executionOOG: this.executionOOG,
      executionGasLimit: this.executionGasLimit,
      userOperationEvent: this.userOperationEvent,
      output: toHex(ctx.output),
      error: ctx.error,
    };
  },

  enter: function enter(frame) {
    if (this._isExecution()) {
      var next = this._depth + 1;
      if (this._executionGasStack[next] === undefined)
        this._executionGasStack[next] = Object.assign({}, this._defaultGasItem);
    }
  },
  exit: function exit(frame) {
    if (this._isExecution()) {
      if (frame.getError() !== undefined) {
        this.reverts.push(toHex(frame.getOutput()));
      }

      if (this._depth >= 2) {
        // Get the final gas item for the nested frame.
        var nested = Object.assign(
          {},
          this._executionGasStack[this._depth + 1] || this._defaultGasItem
        );

        // Reset the nested gas item to prevent double counting on re-entry.
        this._executionGasStack[this._depth + 1] = Object.assign(
          {},
          this._defaultGasItem
        );

        // Keep track of the total gas used by all frames at this depth.
        // This does not account for the gas required due to the 63/64 rule.
        var used = frame.getGasUsed();
        this._executionGasStack[this._depth].used += used;

        // Keep track of the total gas required by all frames at this depth.
        // This accounts for additional gas needed due to the 63/64 rule.
        this._executionGasStack[this._depth].required +=
          used - nested.used + Math.ceil((nested.required * 64) / 63);

        // Keep track of the final gas limit.
        this.executionGasLimit = this._executionGasStack[this._depth].required;
      }
    }
  },

  step: function step(log, db) {
    var opcode = log.op.toString();
    this._depth = log.getDepth();
    if (this._depth === 1 && opcode === "NUMBER") this._marker++;

    if (
      this._depth <= 2 &&
      opcode.startsWith("LOG") &&
      this._isUserOperationEvent(log)
    )
      this._setUserOperationEvent(opcode, log);

    if (log.getGas() < log.getCost() && this._isValidation())
      this.validationOOG = true;

    if (log.getGas() < log.getCost() && this._isExecution())
      this.executionOOG = true;
  },
};
