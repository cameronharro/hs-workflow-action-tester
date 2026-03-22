declare var module: { exports: { main?: (data: any) => any } };
const { exports } = module;

type Envelope = {
  event: Record<string, any>;
  function: string;
};
async function readEnvelope(): Promise<Envelope> {
  const buf = new Uint8Array(1000);
  const read = await Deno.stdin.read(buf);
  if (read === null) {
    throw new Error("No stdin provided");
  }

  const str = new TextDecoder("utf-8").decode(buf.slice(0, read));
  const envelope = JSON.parse(str);
  if (
    typeof envelope === "object" &&
    "event" in envelope &&
    "function" in envelope
  ) {
    return envelope;
  }
  throw new Error(JSON.stringify(`${envelope} is not a valid envelope`));
}

async function main() {
  function callback(data: any) {
    try {
      if (typeof data !== "object" || data === null) {
        throw new Error();
      }
      console.log(JSON.stringify(data));
    } catch {
      throw new Error(
        `expected envelope.function to return JSON data, but received ${data}`,
      );
    }
  }
  const envelope = await readEnvelope();
  eval(envelope.function);
  if (typeof exports.main !== "function") {
    throw new Error("envelope.function did not set exports.main");
  }
  try {
    const result = exports.main(envelope.event);
    callback(result);
  } catch (e) {
    if (e instanceof Error) {
      throw new Error("envelope.function: " + e.message);
    }
    throw e;
  }
  Deno.exit(0);
}

main().catch((e) => {
  console.error(e);
  Deno.exit(1);
});
