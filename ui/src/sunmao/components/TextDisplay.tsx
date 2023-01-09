import { useEffect } from "react";
import { implementRuntimeComponent } from "@sunmao-ui/runtime";
import { Type } from "@sinclair/typebox";

enum Format {
  Plain = "plain",
  Code = "code",
}

export default implementRuntimeComponent({
  version: "custom/v1",
  metadata: {
    name: "TextDisplay",
    displayName: "TextDisplay",
    description: "Display text in different formats.",
    exampleProperties: {
      text: "Test",
      format: Format.Plain,
    },
    annotations: {
      category: "Display",
    },
  },
  spec: {
    properties: Type.Object({
      text: Type.String(),
      format: Type.Enum(Format),
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
    if (format === Format.Code) return <pre>{text}</pre>;

    return <p>{text}</p>;
  };

  return <div ref={elementRef}>{renderContent()}</div>;
});
