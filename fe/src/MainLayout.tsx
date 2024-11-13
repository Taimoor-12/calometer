import { Outlet } from "react-router-dom";
import Navigation from "./components/Navigation";

const MainLayout = () => (
  <div>
    <Navigation />
    <Outlet /> {/* This will render the matched child route component */}
  </div>
);

export default MainLayout;
