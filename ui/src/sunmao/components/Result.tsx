import { implementRuntimeComponent } from "@sunmao-ui/runtime";
import { Type, Static } from "@sinclair/typebox";
import {
  Table as BaseTable,
  Tabs as BaseTabs,
  Tree as BaseTree,
  TreeProps,
} from "@arco-design/web-react";
import { css, cx } from "@emotion/css";
import _ from "lodash";

const BaseTabPane = BaseTabs.TabPane;
const ResultStyle = css`
  .pre-wrapper {
    overflow: auto;
  }
`;

const PropSpec = Type.Object({
  data: Type.Union([
    Type.String(),
    Type.Record(Type.String(), Type.Unknown()),
    Type.Array(Type.Unknown()),
  ]),
});
const StateSpec = Type.Object({
  data: Type.Union([
    Type.String(),
    Type.Record(Type.String(), Type.Unknown()),
    Type.Array(Type.Unknown()),
  ]),
});
const exampleProperties: Static<typeof PropSpec> = {
  data: "",
};

const Result = implementRuntimeComponent({
  version: "cli2ui/v1",
  metadata: {
    name: "result",
    displayName: "Result",
    exampleProperties,
    annotations: {
      category: "Display",
    },
  },
  spec: {
    properties: PropSpec,
    state: {},
    methods: [],
    slots: {},
    styleSlots: ["content"],
    events: [],
  },
})(({ data, elementRef, customStyle }) => {
  const renderContent = (type: "log" | "table" | "tree") => {
    if (type === "log") {
      if (typeof data === "string") return <pre>{data}</pre>;

      try {
        return <pre className="pre-wrapper">{JSON.stringify(data)}</pre>;
      } catch {
        return <p>无法展示为纯文字</p>;
      }
    }
    if (type === "table") {
      if (!Array.isArray(data) || !_.isPlainObject(data[0]))
        return <p>无法展示为表格</p>;

      const columnKeys = Object.keys(data[0]);
      const columns = columnKeys.map((key) => ({
        title: key,
        dataIndex: key,
      }));
      return <BaseTable columns={columns} data={data} />;
    }

    if (type === "tree") {
      if (!_.isPlainObject(data)) return <p>无法展示为树状结构</p>;

      const getTreeData = (
        raw: Record<string, unknown>,
        parentKey?: string
      ): TreeProps["treeData"] => {
        const res: TreeProps["treeData"] = [];
        Object.entries(raw).forEach(([key, value]) => {
          const children = _.isPlainObject(value)
            ? getTreeData(value as Record<string, unknown>, key)
            : [];

          const geTitle = (key: string, value: unknown) => {
            if (Array.isArray(value)) {
              return `${key}: ${value.join(", ")}`;
            }
            if (_.isPlainObject(value)) {
              return key;
            }
            return `${key}: ${value}`;
          };
          const obj = {
            title: geTitle(key, value),
            key: parentKey ? parentKey + "_" + key : key,
            children,
          };
          res.push(obj);
        });
        return res;
      };
      const treeData = getTreeData(data as Record<string, unknown>);
      return <BaseTree treeData={treeData}></BaseTree>;
    }
  };

  return (
    <div
      className={cx(css(customStyle?.content), ResultStyle, "cli2ui-v1-result")}
      ref={elementRef}
    >
      {data ? (
        <BaseTabs>
          <BaseTabPane key="1" title="纯文字">
            {renderContent("log")}
          </BaseTabPane>
          <BaseTabPane key="2" title="表格">
            {renderContent("table")}
          </BaseTabPane>
          <BaseTabPane key="3" title="结构">
            {renderContent("tree")}
          </BaseTabPane>
        </BaseTabs>
      ) : (
        <p>暂无结果</p>
      )}
    </div>
  );
});

export default Result;
