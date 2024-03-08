import { renderToStaticMarkup } from "react-dom/server";
import React from "react";

const App = () => {
  return (
    <div>
      <h1 className="hello jeff">React App</h1>
    </div>
  );
};

const html = renderToStaticMarkup(<App />);
console.log(html);