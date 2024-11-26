import {BrowserRouter as Router, Routes, Route} from 'react-router-dom';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" />
        <Route path="/sign-up"/>
        <Route path="/login" />
        <Route path="/budgets" />
        <Route path="/transactions" />
        <Route path="/forecast" />
      </Routes>
    </Router>
  );
}

export default App;
