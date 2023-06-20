import { Terminal as BaseTerminal } from "xterm";
import { WebLinksAddon } from "xterm-addon-web-links";
import { FitAddon } from "xterm-addon-fit";
import { Type } from "@sinclair/typebox";
import { css, cx } from "@emotion/css";
import { implementRuntimeComponent } from "@sunmao-ui/runtime";

import "xterm/css/xterm.css";
import { useEffect, useRef } from "react";

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

  useEffect(() => {
    terminalRef.current?.reset();
    terminalRef.current?.write(text);
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
    terminal.reset();
    terminal.write(text);

    const fitAddon = new FitAddon();
    terminal.loadAddon(new WebLinksAddon());
    terminal.loadAddon(fitAddon);

    terminal.open(elementRef.current);
    fitAddon.fit();

    terminalRef.current = terminal;
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
