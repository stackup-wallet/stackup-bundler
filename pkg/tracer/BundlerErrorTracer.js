var tracer = {
  reverts: [],
  numberCounter: 0,

  fault: function fault(log, db) {},
  result: function result(ctx, db) {
    return {
      reverts: this.reverts,
    };
  },

  enter: function enter(frame) {},
  exit: function exit(frame) {
    if (frame.getError() !== undefined && this.numberCounter === 3) {
      this.reverts.push(toHex(frame.getOutput()));
    }
  },

  step: function step(log, db) {
    var opcode = log.op.toString();
    if (log.getDepth() === 1 && opcode === "NUMBER") this.numberCounter++;
  },
};
