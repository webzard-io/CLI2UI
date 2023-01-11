import { useEffect } from "react";
import { implementRuntimeComponent, StringUnion } from "@sunmao-ui/runtime";
import { Type } from "@sinclair/typebox";
import { css, cx } from "@emotion/css";

const TextDisplayStyle = css`
  overflow: auto;
`;

export default implementRuntimeComponent({
  version: "custom/v1",
  metadata: {
    name: "TextDisplay",
    displayName: "TextDisplay",
    description: "Display text in different formats.",
    exampleProperties: {
      text: "Test",
      format: "plain",
    },
    annotations: {
      category: "Display",
    },
  },
  spec: {
    properties: Type.Object({
      text: Type.String(),
      format: StringUnion(["plain", "code"]),
    }),
    state: Type.Object({
      text: Type.String(),
    }),
    methods: [],
    slots: {},
    styleSlots: ['content'],
    events: [],
  },
})(({ text, format, mergeState, elementRef, customStyle }) => {
  useEffect(() => {
    mergeState({ text: text });
  }, [mergeState, text]);

  const renderContent = () => {
    if (format === "code") return <pre>{text}</pre>;

    return <p>{text}</p>;
  };

  return (
    <div
      className={cx(css(customStyle?.content), TextDisplayStyle, "custom-text-display")}
      ref={elementRef}
    >
      {renderContent()}
    </div>
  );
});
