// Import the HTTP module
const http = require("http");

// Create a server object
const server = http.createServer((req, res) => {
  let body = "";
  req.on("data", (chunk) => {
    body += chunk.toString();
  });

  req.on("end", () => {
    try {
      console.log(req.method);
      console.log(req.headers);
      console.log(body);
      const obj = JSON.parse(body);
      let responseBody;
      if (obj?.inputFields?.value === "1234") {
        responseBody = JSON.stringify({
          outputFields: { status: "success", hs_execution_state: "SUCCESS" },
        });
      } else {
        responseBody = JSON.stringify({
          outputFields: {
            status: "failure",
            hs_execution_state: "FAIL_CONTINUE",
          },
        });
      }
      res.writeHead(201, {
        "Content-Type": "application/json",
        "Content-Length": Buffer.byteLength(responseBody),
        Connection: "close",
      });
      res.end(responseBody);
    } catch (error) {
      console.error(error);
      res.writeHead(400, { "Content-Type": "application/json" });
      res.end(JSON.stringify({ error: "Invalid JSON" }));
    }
  });
});

// Define the port to listen on const PORT = 9090;
const PORT = 3000;
// Start the server and listen on the specified port
server.listen(PORT, "localhost", () => {
  console.log(`Server running at http://localhost:${PORT}/`);
});
