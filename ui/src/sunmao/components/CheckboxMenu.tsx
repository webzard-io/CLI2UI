import { implementRuntimeComponent } from "@sunmao-ui/runtime";
import { Type, Static } from "@sinclair/typebox";
import { useState, useEffect } from "react";
import {
  Dropdown as BaseDropdown,
  Menu as BaseMenu,
  Checkbox as BaseCheckbox,
  Button as BaseButton,
} from "@arco-design/web-react";
import { IconDown } from "@arco-design/web-react/icon";
import { css, cx } from "@emotion/css";
import _ from "lodash";

const CheckboxMenuStyle = css``;

const PropSpec = Type.Object({
  value: Type.Array(Type.String()),
  text: Type.String(),
  options: Type.Array(
    Type.Object({
      label: Type.String(),
      value: Type.String(),
      disabled: Type.Boolean(),
    })
  ),
});
const StateSpec = Type.Object({
  value: Type.Array(Type.String()),
});

const exampleProperties: Static<typeof PropSpec> = {
  value: ["Item1"],
  text: "Menu",
  options: [
    {
      label: "Item1",
      value: "Item1",
      disabled: false,
    },
  ],
};

const CheckboxMenu = implementRuntimeComponent({
  version: "custom/v1",
  metadata: {
    name: "checkboxMenu",
    displayName: "CheckboxMenu",
    exampleProperties,
    annotations: {
      category: "Data Entry",
    },
  },
  spec: {
    properties: PropSpec,
    state: StateSpec,
    methods: [],
    slots: {},
    styleSlots: ["content"],
    events: [],
  },
})((props) => {
  const { value: defaultValue, options, text, elementRef, mergeState } = props;

  const [value, setValue] = useState<Set<string>>(new Set(defaultValue));

  useEffect(() => {
    mergeState({ value: Array.from(value) });
  }, [mergeState, value]);

  const dropList = (
    <BaseMenu
      onClickMenuItem={() => {
        return false;
      }}
    >
      {options.map((opt) => {
        return (
          <BaseMenu.Item key={opt.value}>
            <BaseCheckbox
              checked={value.has(opt.value)}
              disabled={opt.disabled}
              onClick={() => {
                const newValue = _.clone(value);
                newValue.has(opt.value)
                  ? newValue.delete(opt.value)
                  : newValue.add(opt.value);
                setValue(newValue);
              }}
            >
              {opt.label}
            </BaseCheckbox>
          </BaseMenu.Item>
        );
      })}
    </BaseMenu>
  );

  return (
    <div
      ref={elementRef}
      className={cx(CheckboxMenuStyle, "custom-v1-check-box")}
    >
      <BaseDropdown droplist={dropList} position="br">
        <BaseButton type="text">
          {text} <IconDown />
        </BaseButton>
      </BaseDropdown>
    </div>
  );
});

export default CheckboxMenu;
