import { useState } from "react";

function App() {
  const [count, setCount] = useState(0);

  return (
    <div>
      <h1 className="text-3xl font-bold underline">App Store RSS Reviews</h1>
      <button onClick={() => setCount(count + 1)}>Count</button>
      <p>Count: {count}</p>
    </div>
  );
}

export default App;
