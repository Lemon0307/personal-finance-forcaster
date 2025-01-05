import {BrowserRouter as Router, Routes, Route} from 'react-router-dom';
import {Home, Forecast, Budgets, Login, SignUp, Transactions, AddBudget, UpdateBudget} from './pages';
import {Navbar} from './components';

function App() {
  return (
    <div className="p-10">
      <Router>
        <Navbar />
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/sign-up" element={<SignUp />} />
          <Route path="/login" element={<Login />} />
          <Route path="/budgets" element={<Budgets />} />
          <Route path="/budgets/add-budget" element={<AddBudget />} />
          <Route path="/budgets/update-budget" element={<UpdateBudget />}/>
          <Route path="/transactions" element={<Transactions />} />
          <Route path="/forecast" element={<Forecast />} />
        </Routes>
      </Router>
    </div>
  );
}

export default App;
