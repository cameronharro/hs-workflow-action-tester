const module = { exports: {} };
const { exports } = module;

type Envelope = {
  event: Record<string, any>;
  function: string;
};
async function readEnvelope(): Promise<Envelope> {
  const buf = new Uint8Array(100);
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
  try {
    const envelope = await readEnvelope();
    eval(envelope.function);
    Deno.exit(0);
  } catch (e) {
    console.error(e);
    Deno.exit(1);
  }
}

await main();
