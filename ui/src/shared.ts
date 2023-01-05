import { sunmaoChakraUILib } from "@sunmao-ui/chakra-ui-lib";
import { ArcoDesignLib } from "@sunmao-ui/arco-lib";
import type { Application, Module } from "@sunmao-ui/core";
import {
  implementUtilMethod,
  initSunmaoUI,
  UtilMethodFactory,
} from "@sunmao-ui/runtime";
import { useEffect } from "react";
import * as jdp from "jsondiffpatch";
import CustomLib from "./sunmao/lib";

export function getLibs({
  ws,
  handlers,
  utilMethods,
}: {
  ws: WebSocket;
  handlers: string[];
  utilMethods?: UtilMethodFactory[];
}) {
  return [
    sunmaoChakraUILib,
    ArcoDesignLib,
    CustomLib,
    {
      utilMethods: (utilMethods || []).concat(
        handlers.map(
          (handler) => () =>
            implementUtilMethod({
              version: "binding/v1",
              metadata: {
                name: handler,
              },
              spec: {
                parameters: {} as any,
              },
            })((params) => {
              ws.send(
                JSON.stringify({
                  type: "Action",
                  handler,
                  params,
                })
              );
            })
        )
      ),
    },
  ];
}

export type ServerMessage = {
  type: string;
  componentId: string;
  name: string;
  parameters?: any;
};

export function useApiService({
  ws,
  apiService,
}: {
  ws: WebSocket;
  apiService: ReturnType<typeof initSunmaoUI>["apiService"];
}) {
  useEffect(() => {
    const messageHandler = (evt: MessageEvent) => {
      try {
        const message: ServerMessage = JSON.parse(evt.data);
        if (message.type !== "UiMethod") {
          return;
        }
        apiService.send("uiMethod", {
          componentId: message.componentId,
          name: message.name,
          parameters: message.parameters,
        });
      } catch (error) {
        console.log("message handler", error);
      }
    };
    ws.addEventListener("message", messageHandler);
    return () => ws.removeEventListener("message", messageHandler);
  }, [apiService]);
}

export type BaseProps = {
  handlers: string[];
  ws: WebSocket;
  utilMethods?: UtilMethodFactory[];
} & Pick<
  MainOptions,
  "application" | "modules" | "applicationPatch" | "modulesPatch"
>;

export type MainOptions = {
  application: any;
  modules: any[];
  wsUrl: string;
  reloadWhenWsDisconnected: boolean;
  handlers: string[];
  utilMethods?: { options: any; impl: any }[];
  applicationPatch?: any;
  modulesPatch?: any;
};

// const PREFIX = "/sunmao-binding-patch";
const PREFIX = "/sunmao-fs";
export const APPLICATION_NAME = 'app';

const diffpatcher = jdp.create({
  objectHash: function (obj: any, index: number) {
    if (typeof obj._id !== "undefined") {
      return obj._id;
    }
    if (typeof obj.id !== "undefined") {
      return obj.id;
    }
    if (typeof obj.name !== "undefined") {
      return obj.name;
    }
    return "$$index:" + index;
  },
  cloneDiffValues: true,
});

export function saveApp(app: Application, base: Application) {
  return fetch(`${PREFIX}/${APPLICATION_NAME}`, {
    method: "put",
    headers: {
      "content-type": "application/json",
    },
    body: JSON.stringify({
      // delta: diffpatcher.diff(base, app),
      value: app,
    }),
  });
}

export function saveModules(modules: Module[], base: Module[]) {
  return fetch(`${PREFIX}/modules`, {
    method: "put",
    headers: {
      "content-type": "application/json",
    },
    body: JSON.stringify({
      // delta: diffpatcher.diff(base, modules),
      value: modules,
    }),
  });
}

function isEmptyDelta(delta?: jdp.Delta) {
  return !delta || Object.keys(delta).length === 0;
}

export function patchApp(base: Application, delta?: jdp.Delta): Application {
  return isEmptyDelta(delta)
    ? base
    : diffpatcher.patch(diffpatcher.clone(base), delta!);
}

export function patchModules(base: Module[], delta?: jdp.Delta): Module[] {
  return isEmptyDelta(delta)
    ? base
    : diffpatcher.patch(diffpatcher.clone(base), delta!);
}

/**
 * get existed sunmao application schema
 * @param name 
 * @returns 
 */
export async function fetchApp(name: string): Promise<Application> {
  const application = await (await fetch(`${PREFIX}/${name}`)).json();

  if (application.kind === 'Application') {
    return application;
  }

  throw new Error('failed to load schema');
}

/**
 * get existed sunmao module schema
 * @param name 
 * @returns 
 */
export async function fetchModules(): Promise<Module[]> {
  const response = await (await fetch(`${PREFIX}/modules`)).json();

  if (Array.isArray(response)) {
    return response;
  }

  throw new Error('failed to load schema');
}
