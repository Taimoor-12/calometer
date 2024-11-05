import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import "./App.css";
import Signup from "./Signup";
import Login from "./Login";
import { ToastContainer } from "react-toastify";
import 'react-toastify/dist/ReactToastify.css';
import AddBodyDetails from "./AddBodyDetails";

function App() {
  return (
    <Router>
      <ToastContainer position="bottom-right" autoClose={3000} />
      <Routes>
        <Route path="/signup" element={<Signup/>} />
        <Route path="/login" element={<Login/>} />
        <Route path="/addBodyDetails" element={<AddBodyDetails/>} />
        <Route path="/dashboard" element={<><p>Hello World</p></>} />
      </Routes>
    </Router>
  );
}

export default App;
