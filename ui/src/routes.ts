import { APPLICATION_NAME } from "./shared";
export type RouteConfig = {
  name: string;
  path: string;
};

export type RedirectConfig = {
  path: string;
  redirect: string;
};

export type Config = RouteConfig | RedirectConfig;

const routes: Config[] = [
  {
    path: "*",
    name: APPLICATION_NAME,
  },
];

export default routes;
