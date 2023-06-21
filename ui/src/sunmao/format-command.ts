type PlainValue = string | number | boolean;

type FlagOrArg = {
  value: PlainValue | PlainValue[];
  isArg: boolean;
  name: string;
  short?: string;
};

export function formatCommand(
  cmdName: string,
  values: Record<string, FlagOrArg>,
  checked: string[]
) {
  const cmd = [cmdName];

  for (const key in values) {
    if (!checked.includes(key)) {
      continue;
    }

    const v = values[key];
    if (Array.isArray(v.value)) {
      for (const childValue of v.value) {
        if (v.isArg) {
          cmd.push(getKey(v));
        }
        cmd.push(String(childValue));
      }
    } else {
      if (v.isArg) {
        cmd.push(getKey(v));
      }
      cmd.push(String(v.value));
    }
  }

  // console.log({ cmdName, values, checked }, cmd.join(" "));
  return cmd.join(" ").trim();
}

function getKey(foa: FlagOrArg) {
  const key = foa.short ? `-${foa.short}` : `--${foa.name}`;
  return key;
}
