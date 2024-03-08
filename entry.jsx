import { renderToString } from "react-dom/server.browser";

import React from "react";

const App = () => {
  return (
    <div>
      <h1 className="hello jeff">React App</h1>
    </div>
  );
};

renderToString(<App />);