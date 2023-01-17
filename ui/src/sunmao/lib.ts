import { type SunmaoLib } from "@sunmao-ui/runtime";
import TextDisplay from "./components/TextDisplay";
import ArrayInput from "./components/ArrayInput";
import CheckboxMenu from "./components/CheckboxMenu";
import "./style.css";

const lib: SunmaoLib = {
  components: [TextDisplay, ArrayInput, CheckboxMenu],
  traits: [],
};

export default lib;
