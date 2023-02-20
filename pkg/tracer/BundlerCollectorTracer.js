// This is the same BundlerCollectorTracer from github.com/eth-infinitism/bundler transpiled down to ES5.

var tracer = {
  numberLevels: [],
  currentLevel: null,
  keccak: [],
  calls: [],
  logs: [],
  debug: [],
  lastOp: "",
  numberCounter: 0,

  fault: function fault(log, db) {
    this.debug.push(
      "fault depth=" +
        log.getDepth() +
        " gas=" +
        log.getGas() +
        " cost=" +
        log.getCost() +
        " err=" +
        log.getError()
    );
  },
  result: function result(ctx, db) {
    return {
      numberLevels: this.numberLevels,
      keccak: this.keccak,
      logs: this.logs,
      calls: this.calls,
      // for internal debugging.
      debug: this.debug,
    };
  },
  enter: function enter(frame) {
    this.debug.push(
      "enter gas=" +
        frame.getGas() +
        " type=" +
        frame.getType() +
        " to=" +
        toHex(frame.getTo()) +
        " in=" +
        toHex(frame.getInput()).slice(0, 500)
    );
    this.calls.push({
      type: frame.getType(),
      from: toHex(frame.getFrom()),
      to: toHex(frame.getTo()),
      method: toHex(frame.getInput()).slice(0, 10),
      gas: frame.getGas(),
      value: frame.getValue(),
    });
  },
  exit: function exit(frame) {
    this.calls.push({
      type: frame.getError() != null ? "REVERT" : "RETURN",
      gasUsed: frame.getGasUsed(),
      data: toHex(frame.getOutput()).slice(0, 1000),
    });
  },

  // increment the "key" in the list. if the key is not defined yet, then set it to "1"
  countSlot: function countSlot(list, key) {
    if (!list[key]) list[key] = 0;
    list[key] += 1;
  },
  step: function step(log, db) {
    var opcode = log.op.toString();
    // this.debug.push(this.lastOp + '-' + opcode + '-' + log.getDepth() + '-' + log.getGas() + '-' + log.getCost())
    if (log.getGas() < log.getCost()) {
      this.currentLevel.oog = true;
    }

    if (opcode === "REVERT" || opcode === "RETURN") {
      if (log.getDepth() === 1) {
        // exit() is not called on top-level return/revert, so we reconstruct it
        // from opcode
        var ofs = parseInt(log.stack.peek(0).toString());
        var len = parseInt(log.stack.peek(1).toString());
        var data = toHex(log.memory.slice(ofs, ofs + len)).slice(0, 1000);
        this.debug.push(opcode + " " + data);
        this.calls.push({
          type: opcode,
          gasUsed: 0,
          data: data,
        });
      }
    }

    if (opcode.match(/^(EXT.*|CALL|CALLCODE|DELEGATECALL|STATICCALL)$/) != null) {
      // this.debug.push('op=' + opcode + ' last=' + this.lastOp + ' stacksize=' + log.stack.length())
      var idx = opcode.startsWith("EXT") ? 0 : 1;
      var addr = toAddress(log.stack.peek(idx).toString(16));
      var addrHex = toHex(addr);
      var contractSize =
        (contractSize = this.currentLevel.contractSize[addrHex]) !== null &&
        contractSize !== void 0
          ? contractSize
          : 0;
      if (contractSize === 0 && !isPrecompiled(addr)) {
        this.currentLevel.contractSize[addrHex] = db.getCode(addr).length;
      }
    }

    if (log.getDepth() === 1) {
      // NUMBER opcode at top level split levels
      if (opcode === "NUMBER") this.numberCounter++;
      if (this.numberLevels[this.numberCounter] == null) {
        this.currentLevel = this.numberLevels[this.numberCounter] = {
          access: {},
          opcodes: {},
          contractSize: {},
        };
      }
      this.lastOp = "";
      return;
    }

    if (this.lastOp === "GAS" && !opcode.includes("CALL")) {
      // count "GAS" opcode only if not followed by "CALL"
      this.countSlot(this.currentLevel.opcodes, "GAS");
    }
    if (opcode !== "GAS") {
      // ignore "unimportant" opcodes:
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
      var _addr = toHex(log.contract.getAddress());
      var access = void 0;
      if ((access = this.currentLevel.access[_addr]) == null) {
        this.currentLevel.access[_addr] = access = {
          reads: {},
          writes: {},
        };
      }
      this.countSlot(opcode === "SLOAD" ? access.reads : access.writes, slot);
    }

    if (opcode === "KECCAK256") {
      // collect keccak on 64-byte blocks
      var _ofs = parseInt(log.stack.peek(0).toString());
      var _len = parseInt(log.stack.peek(1).toString());
      // currently, solidity uses only 2-word (6-byte) for a key. this might change..
      // still, no need to return too much
      if (_len > 20 && _len < 512) {
        // if (len == 64) {
        this.keccak.push(toHex(log.memory.slice(_ofs, _ofs + _len)));
      }
    } else if (opcode.startsWith("LOG")) {
      var count = parseInt(opcode.substring(3));
      var _ofs2 = parseInt(log.stack.peek(0).toString());
      var _len2 = parseInt(log.stack.peek(1).toString());
      var topics = [];
      for (var i = 0; i < count; i++) {
        // eslint-disable-next-line @typescript-eslint/restrict-plus-operands
        topics.push("0x" + log.stack.peek(2 + i).toString(16));
      }
      var _data = toHex(log.memory.slice(_ofs2, _ofs2 + _len2));
      this.logs.push({
        topics: topics,
        data: _data,
      });
    }
  },
};
