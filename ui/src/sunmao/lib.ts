import { type SunmaoLib } from "@sunmao-ui/runtime";
import TextDisplay from "./components/TextDisplay";
import ArrayInput from "./components/ArrayInput";
import CheckboxMenu from "./components/CheckboxMenu";
import Result from "./components/Result";
import "./style.css";

const lib: SunmaoLib = {
  components: [TextDisplay, ArrayInput, CheckboxMenu, Result],
  traits: [],
};

export default lib;
