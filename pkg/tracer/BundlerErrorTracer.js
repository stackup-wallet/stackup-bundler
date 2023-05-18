var tracer = {
  reverts: [],
  validationOOG: false,
  executionOOG: false,

  _marker: 0,
  _validationMarker: 1,
  _executionMarker: 3,

  _isValidation: function () {
    return (
      this._marker >= this._validationMarker &&
      this._marker < this._executionMarker
    );
  },

  _isExecution: function () {
    return this._marker === this._executionMarker;
  },

  fault: function fault(log, db) {},
  result: function result(ctx, db) {
    return {
      reverts: this.reverts,
      validationOOG: this.validationOOG,
      executionOOG: this.executionOOG,
      output: toHex(ctx.output),
    };
  },

  enter: function enter(frame) {},
  exit: function exit(frame) {
    if (frame.getError() !== undefined && this._isExecution()) {
      this.reverts.push(toHex(frame.getOutput()));
    }
  },

  step: function step(log, db) {
    var opcode = log.op.toString();
    if (log.getDepth() === 1 && opcode === "NUMBER") this._marker++;

    if (log.getGas() < log.getCost() && this._isValidation())
      this.validationOOG = true;

    if (log.getGas() < log.getCost() && this._isExecution())
      this.executionOOG = true;
  },
};
