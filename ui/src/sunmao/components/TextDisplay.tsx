import { useEffect } from "react";
import { implementRuntimeComponent, StringUnion } from "@sunmao-ui/runtime";
import { Type } from "@sinclair/typebox";

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
    styleSlots: [],
    events: [],
  },
})(({ text, format, mergeState, elementRef }) => {
  useEffect(() => {
    mergeState({ text: text });
  }, [mergeState, text]);

  const renderContent = () => {
    if (format === 'code') return <pre>{text}</pre>;

    return <p>{text}</p>;
  };

  return <div ref={elementRef}>{renderContent()}</div>;
});
