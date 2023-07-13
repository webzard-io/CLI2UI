import { Terminal as BaseTerminal } from "xterm";
import { WebLinksAddon } from "xterm-addon-web-links";
import { FitAddon } from "xterm-addon-fit";
import { Type } from "@sinclair/typebox";
import { css, cx } from "@emotion/css";
import { implementRuntimeComponent } from "@sunmao-ui/runtime";

import "xterm/css/xterm.css";
import { useEffect, useRef } from "react";

class SequentialWriter {
  private currentWrite: Promise<void> | null = null;
  private pendingText: string | null = null;
  private terminal: BaseTerminal;

  constructor(terminal: BaseTerminal) {
    this.terminal = terminal;
  }

  public write(text: string) {
    // when there is no active writer, start writing
    if (!this.currentWrite) {
      this.currentWrite = new Promise((resolve) => {
        this.terminal.reset();
        this.terminal.write(text, () => {
          resolve();
          // reset active writer
          this.currentWrite = null;

          // when there is some pending text, write it and reset the pending state
          if (this.pendingText) {
            this.write(this.pendingText);
            this.pendingText = null;
          }
        });
      });
      return;
    }

    /**
     * Keep the text as pending state so we can handle it later.
     *
     * During the current active write process, we will only
     * keep the latest pending text and drop the expired texts.
     */
    this.pendingText = text;
  }
}

export default implementRuntimeComponent({
  version: "cli2ui/v1",
  metadata: {
    name: "terminal",
    displayName: "Terminal",
    description: "Display result in a terminal.",
    exampleProperties: {
      text: "Hello CLI2UI",
    },
    annotations: {
      category: "Display",
    },
  },
  spec: {
    properties: Type.Object({
      text: Type.String(),
    }),
    state: Type.Object({
      text: Type.String(),
    }),
    methods: [],
    slots: {},
    styleSlots: ["content"],
    events: [],
  },
})(({ text = "", mergeState, elementRef, customStyle }) => {
  const terminalRef = useRef<BaseTerminal | null>(null);
  const writerRef = useRef<SequentialWriter | null>(null);

  useEffect(() => {
    writerRef.current?.write(text);
    mergeState({ text: text });
  }, [mergeState, text]);

  useEffect(() => {
    // init terminal

    if (!elementRef?.current) {
      return;
    }
    if (terminalRef.current) {
      return;
    }

    const terminal = new BaseTerminal({
      // cursorBlink: true,
      convertEol: true,
      // cols: 80,
      cursorWidth: 1,
    });

    const fitAddon = new FitAddon();
    terminal.loadAddon(new WebLinksAddon());
    terminal.loadAddon(fitAddon);

    terminal.open(elementRef.current);
    fitAddon.fit();

    terminalRef.current = terminal;
    writerRef.current = new SequentialWriter(terminal);

    writerRef.current.write(text);
  }, [text]);

  return (
    <div
      id="terminal"
      ref={elementRef}
      className={cx(
        css`
          width: 100%;
          height: 100%;

          .xterm .terminal-cursor {
            opacity: 0;
          }
        `,
        css(customStyle?.content)
      )}
    />
  );
});
