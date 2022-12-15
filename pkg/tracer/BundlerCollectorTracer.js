// This is the same BundlerCollectorTracer from github.com/eth-infinitism/bundler transpiled down to ES5.

var tracer = {
  numberLevels: {},
  currentLevel: null,
  keccak: [],
  calls: [],
  logs: [],
  debug: [],
  lastOp: "",
  numberCounter: 0,
  count: 0,

  fault: function (log, db) {
    this.debug.push(["fault", log.getError()]);
  },

  result: function (ctx, db) {
    return {
      numberLevels: this.numberLevels,
      keccak: this.keccak,
      logs: this.logs,
      calls: this.calls,
      debug: this.debug,
    };
  },

  enter: function (frame) {
    this.debug.push([
      "enter " +
        frame.getType() +
        " " +
        toHex(frame.getTo()) +
        " " +
        toHex(frame.getInput()).slice(0, 100),
    ]);

    this.calls.push({
      type: frame.getType(),
      from: toHex(frame.getFrom()),
      to: toHex(frame.getTo()),
      value: frame.getValue(),
    });
  },

  exit: function (frame) {
    this.debug.push(
      "exit err=" + frame.getError() + ", gas=" + frame.getGasUsed()
    );
  },

  countSlot: function countSlot(list, key) {
    list[key] += typeof list[key] == "number" ? 0 : 1;
  },

  step: function step(log, db) {
    var opcode = log.op.toString();

    if (opcode === "NUMBER") this.numberCounter++;

    if (this.numberLevels[this.numberCounter] == null) {
      this.currentLevel = this.numberLevels[this.numberCounter] = {
        access: {},
        opcodes: {},
      };
    }

    if (log.getDepth() === 1) {
      return;
    }

    if (this.lastOp === "GAS" && !opcode.includes("CALL")) {
      this.countSlot(this.currentLevel.opcodes, "GAS");
    }

    if (opcode !== "GAS") {
      if (
        opcode.match(
          /^(DUP\d+|PUSH\d+|SWAP\d+|POP|ADD|SUB|MUL|DIV|EQ|LTE?|S?GTE?|SLT|SH[LR]|AND|OR|NOT|ISZERO)$/
        ) == null
      ) {
        this.countSlot(this.currentLevel.opcodes, opcode);
      }
    }

    this.lastOp = opcode;

    if (opcode === "SLOAD" || opcode === "SSTORE") {
      var slot = log.stack.peek(0).toString(16);
      var addr = toHex(log.contract.getAddress());
      var access = void 0;

      if ((access = this.currentLevel.access[addr]) == null) {
        this.currentLevel.access[addr] = access = {
          reads: {},
          writes: {},
        };
      }

      this.countSlot(opcode === "SLOAD" ? access.reads : access.writes, slot);
    }

    if (opcode === "REVERT" || opcode === "RETURN") {
      var ofs = parseInt(log.stack.peek(0).toString());
      var len = parseInt(log.stack.peek(1).toString());

      this.debug.push(
        opcode + " " + toHex(log.memory.slice(ofs, ofs + len)).slice(0, 100)
      );
    } else if (opcode === "KECCAK256") {
      var _ofs = parseInt(log.stack.peek(0).toString());
      var _len = parseInt(log.stack.peek(1).toString());

      if (_len < 512) {
        this.keccak.push(toHex(log.memory.slice(_ofs, _ofs + _len)));
      }
    } else if (opcode.startsWith("LOG")) {
      var count = parseInt(opcode.substring(3));
      var _ofs2 = parseInt(log.stack.peek(0).toString());
      var _len2 = parseInt(log.stack.peek(1).toString());

      var topics = [];
      for (var i = 0; i < count; i++) {
        topics.push("0x" + log.stack.peek(2 + i).toString(16));
      }

      this.logs.push({
        topics: topics,
        data: toHex(log.memory.slice(_ofs2, _ofs2 + _len2)),
      });
    }
  },
};
