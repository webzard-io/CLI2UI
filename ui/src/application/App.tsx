import { initSunmaoUI } from "@sunmao-ui/runtime";
import "@sunmao-ui/arco-lib/dist/index.css";
import {
  getLibs,
  useApiService,
  BaseProps,
  // patchApp,
  patchModules,
  fetchApp,
  APPLICATION_NAME,
} from "../shared";
import { RuntimeModule } from "@sunmao-ui/core";
import { type Application } from "@sunmao-ui/core";
import { useState, useEffect } from "react";

function App(props: BaseProps) {
  const {
    application,
    modules,
    handlers,
    ws,
    utilMethods,
    // applicationPatch,
    modulesPatch,
  } = props;

  const [_app, setApp] = useState<Application>(application);
  useEffect(() => {
    (async function () {
      const [app] = await Promise.all([fetchApp(APPLICATION_NAME)]);

      setApp(app);
    })();
  }, []);

  const {
    App: SunmaoApp,
    apiService,
    registry,
  } = initSunmaoUI({
    libs: getLibs({ ws, handlers, utilMethods }),
  });

  if (modules) {
    patchModules(modules, modulesPatch).forEach((moduleSchema) => {
      registry.registerModule(moduleSchema as RuntimeModule);
    });
  }

  useApiService({ ws, apiService });

  return <SunmaoApp options={_app} />;
}

export default App;
