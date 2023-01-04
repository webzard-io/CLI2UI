import React from "react";
import {
  implementRuntimeComponent,
  PRESET_PROPERTY_CATEGORY,
} from "@sunmao-ui/runtime";
import { Type } from "@sinclair/typebox";

export default implementRuntimeComponent({
  version: "custom/v1",
  metadata: {
    name: "test",
    displayName: "Test",
    description: "A test component.",
    exampleProperties: {
      text: "Test",
    },
    annotations: {
      category: PRESET_PROPERTY_CATEGORY.Basic,
    },
  },
  spec: {
    properties: Type.Object({
      text: Type.String({ title: "Text" }),
    }),
    state: Type.Object({
      text: Type.String(),
    }),
    methods: [],
    slots: {},
    styleSlots: [],
    events: [],
  },
})(({ text }) => {
  return <div>{text}</div>;
});
