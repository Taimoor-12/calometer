import { useState, useEffect, useCallback } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import { http_get, http_post } from "./lib/http";
import { toast } from "react-toastify";
import Navigation from "./components/Navigation";
import s from "./Dashboard.module.css";
import CreateLogModal from "./components/CreateLogModal";

interface UserCalorieLog {
  LogDate: string;
  CaloriesBurnt: number;
  CaloriesConsumed: number;
  Tdee: number;
  Updated_at: string;
  LogStatus: string;
}

interface GetCaloricLogsHandlerResp {
  monthly_logs: {
    [monthYear: string]: UserCalorieLog[];
  };
}

const Dashboard = () => {
  const apiUrl = process.env.REACT_APP_API_URL;
  const navigate = useNavigate();
  const location = useLocation();

  const [netCaloricBalance, setNetCaloricBalance] = useState(0);
  const [calorieLogs, setCalorieLogs] = useState<GetCaloricLogsHandlerResp>();
  const [isAddLogModalOpen, setIsAddLogModalOpen] = useState(false);

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

  const getCalorieLogsCall = useCallback(async () => {
    try {
      const resp = await http_get(`${apiUrl}/api/users/log/get`)
      const respCode = +Object.keys(resp.code)[0]
      if (respCode === 200) {
        console.log(resp);
        setCalorieLogs(resp.data.monthly_logs)
      }
    } catch (e) {
      toast.error("Something went wrong, please try again.");
    }
  }, [apiUrl])

  useEffect(() => {
    netCaloricBalanceCall();
  }, [netCaloricBalanceCall]);

  useEffect(() => {
    getCalorieLogsCall();
  }, [getCalorieLogsCall])

  useEffect(() => {
    if (location.state?.from === "login") {
      return;
    }

    const loginUser = async () => {
      try {
        const resp = await http_post(`${apiUrl}/api/users/login`, {});
        const respCode = +Object.keys(resp.code)[0];
        if (respCode !== 200) {
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

  const toggleModalOpen = () => {
    setIsAddLogModalOpen((prev) => !prev)
  }

  const months = calorieLogs ? Object.keys(calorieLogs) : []

  return (
    <div className={s.main}>
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
      <div className={ months.length !== 0  ? s.tilesWrapperDiv : s.noLogsWrapper}>
        {months.length !== 0 ? months.map((month) => (
          <div key={month} className={s.tilesDiv}>
            <p>{month}</p>
          </div>
        )) : <p>No logs exist</p>}
      </div>
      <div className={s.logAddDiv} onClick={toggleModalOpen}>
        <button>+</button>
      </div>

      {isAddLogModalOpen ? <CreateLogModal toggleModalOpen={toggleModalOpen}/> : null}
    </div>
  );
};

export default Dashboard;
