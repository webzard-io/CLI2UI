import { implementRuntimeComponent, StringUnion } from "@sunmao-ui/runtime";
import { Type, Static } from "@sinclair/typebox";
import { useState, useCallback, useEffect } from "react";
import {
  Button as BaseButton,
  Input as BaseInput,
  InputNumber as BaseInputNumber,
  Space,
} from "@arco-design/web-react";
import { IconClose } from "@arco-design/web-react/icon";
import { css, cx } from "@emotion/css";

const ArrayInputStyle = css`
  .inputs-wrapper {
    margin-bottom: 8px;
  }
  .input-wrapper {
    display: flex;
    align-items: center;
    gap: 8px;
  }
  .remove-btn {
    cursor: pointer;
  }
  .arco-space {
    width: 100%;
  }
`;

const PropSpec = Type.Object({
  value: Type.Union([Type.Array(Type.String()), Type.Array(Type.Integer())]),
  type: StringUnion(["number", "string"]),
  placeholder: Type.String(),
  disabled: Type.Boolean(),
});
const StateSpec = Type.Object({
  value: Type.Union([Type.Array(Type.String()), Type.Array(Type.Integer)]),
});

const exampleProperties: Static<typeof PropSpec> = {
  value: [""],
  type: "string",
  placeholder: "Please input",
  disabled: false,
};

const ArrayInput = implementRuntimeComponent({
  version: "custom/v1",
  metadata: {
    name: "arrayInput",
    displayName: "ArrayInput",
    exampleProperties,
    annotations: {
      category: "Data Entry",
    },
  },
  spec: {
    properties: PropSpec,
    state: StateSpec,
    methods: {
      setValue: Type.Object({
        value: Type.Union([
          Type.Array(Type.String()),
          Type.Array(Type.Integer),
        ]),
      }),
    },
    slots: {},
    styleSlots: ["input"],
    events: ["onChange", "onBlur", "onAdd", "onRemove"],
  },
})((props) => {
  const {
    value: defaultValue,
    type,
    disabled,
    placeholder,
    elementRef,
    mergeState,
    callbackMap,
  } = props;

  const [value, setValue] = useState<number[] | string[]>(defaultValue);

  const onAdd = useCallback(() => {
    if (type === "number") {
      const newValue = [...(value as number[])];
      newValue.push(0);
      setValue(newValue);
      // @ts-ignore
      mergeState({ value: newValue });
    } else {
      const newValue = [...(value as string[])];
      newValue.push("");
      setValue(newValue);
      mergeState({ value: newValue });
    }

    if (callbackMap?.onAdd) {
      callbackMap.onAdd();
    }
  }, [value, setValue, mergeState]);

  const onRemove = useCallback(
    (idx) => {
      const newValue = [...value] as string[] | number[];
      newValue.splice(idx, 1);
      setValue(newValue);
      // @ts-ignore
      mergeState({ value: newValue });
      if (callbackMap?.onRemove) {
        callbackMap.onRemove();
      }
    },
    [value, setValue, mergeState]
  );

  const onChange = useCallback(
    (idx, val) => {
      const newValue = [...value] as string[] | number[];
      newValue.splice(idx, 1, val);
      setValue(newValue);
      // @ts-ignore
      mergeState({ value: newValue });
      if (callbackMap?.onChange) {
        callbackMap.onChange();
      }
    },
    [value, setValue, mergeState]
  );

  const renderStringInputs = () => {
    return (value as string[]).map((item, idx) => {
      return (
        <div className="input-wrapper" key={idx}>
          <BaseInput
            value={item}
            placeholder={placeholder}
            disabled={disabled}
            onChange={(val) => onChange(idx, val)}
          ></BaseInput>
          {value.length > 1 ? (
            <IconClose
              className="remove-btn"
              onClick={() => onRemove(idx)}
            ></IconClose>
          ) : null}
        </div>
      );
    });
  };

  const renderNumberInputs = () => {
    return (value as number[]).map((item, idx) => {
      return (
        <div className="input-wrapper" key={idx}>
          <BaseInputNumber
            value={item}
            placeholder={placeholder}
            disabled={disabled}
            onChange={(val) => onChange(idx, val)}
          ></BaseInputNumber>
          {value.length > 1 ? (
            <IconClose
              className="remove-btn"
              onClick={() => onRemove(idx)}
            ></IconClose>
          ) : null}
        </div>
      );
    });
  };

  return (
    <div
      ref={elementRef}
      className={cx(ArrayInputStyle, "custom-v1-array-input")}
    >
      <div className="inputs-wrapper">
        <Space direction="vertical">
          {type === "number" ? renderNumberInputs() : renderStringInputs()}
        </Space>
      </div>

      <BaseButton className="add-input-btn" disabled={disabled} onClick={onAdd}>
        添加
      </BaseButton>
    </div>
  );
});

export default ArrayInput;
