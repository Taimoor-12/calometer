import { useState, useEffect, useCallback } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import { http_get, http_post } from "./lib/http";
import { toast } from "react-toastify";
import Navigation from "./components/Navigation";
import s from "./Dashboard.module.css";

const Dashboard = () => {
  const apiUrl = process.env.REACT_APP_API_URL;
  const navigate = useNavigate();
  const location = useLocation();

  const [netCaloricBalance, setNetCaloricBalance] = useState(0);

  const netCaloricBalanceCall = useCallback(async () => {
    try {
      const resp = await http_get(
        `${apiUrl}/api/users/net_caloric_balance/get`
      );
      const respCode = +Object.keys(resp.code)[0];
      if (respCode === 200) {
        setNetCaloricBalance(resp.data.net_caloric_balance);
      }
    } catch (e) {
      toast.error("Something went wrong, please try again.");
    }
  }, [apiUrl]);

  useEffect(() => {
    netCaloricBalanceCall();
  }, [netCaloricBalanceCall]);

  useEffect(() => {
    if (location.state?.from === "login") {
      return;
    }

    const loginUser = async () => {
      try {
        const resp = await http_post(`${apiUrl}/api/users/login`, {});
        const respCode = +Object.keys(resp.code)[0];
        if (respCode === 401) {
          navigate("/login", {
            state: { from: "dashboard" },
          });
        }
      } catch (error) {
        console.error("Login failed:", error);
        toast.error("Something went wrong, please try again.");
      }
    };

    loginUser();
  }, [location.state, apiUrl, navigate]);

  return (
    <div className={s.main}>
      <Navigation />
      <div className={s.netBalanceDiv}>
        <p>
          NET CALORIC BALANCE:{" "}
          <span
            className={
              netCaloricBalance > 0
                ? s.plus
                : netCaloricBalance < 0
                ? s.minus
                : ""
            }
          >
            {netCaloricBalance > 0 ? "+" : netCaloricBalance < 0 ? "-" : ""}
          </span>
          {netCaloricBalance}
        </p>
      </div>
      <div className={s.tilesWrapperDiv}>
        <div className={s.tilesDiv}>
          <p>September, 2024</p>
        </div>
        <div className={s.tilesDiv}>
          <p>October, 2024</p>
        </div>
        <div className={s.tilesDiv}>
          <p>November, 2024</p>
        </div>
        <div className={s.tilesDiv}>
          <p>November, 2024</p>
        </div>
      </div>
      <div className={s.logAddDiv}>
        <button>+</button>
      </div>
    </div>
  );
};

export default Dashboard;
