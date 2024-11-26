import { useState, useCallback, useEffect } from "react";
import { http_get } from "../lib/http";
import { toast } from "react-toastify";
import s from "./NetCaloricBalance.module.css";

const NetCaloricBalance = () => {
  const apiUrl = process.env.REACT_APP_API_URL;

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

  return (
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
  );
}

export default NetCaloricBalance