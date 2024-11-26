import {BrowserRouter as Router, Routes, Route} from 'react-router-dom';
import {Home, Forecast, Budgets, Login, SignUp, Transactions} from './pages';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/sign-up" element={<SignUp />} />
        <Route path="/login" element={<Login />} />
        <Route path="/budgets" element={<Budgets />} />
        <Route path="/transactions" element={<Transactions />} />
        <Route path="/forecast" element={<Forecast />} />
      </Routes>
    </Router>
  );
}

export default App;
