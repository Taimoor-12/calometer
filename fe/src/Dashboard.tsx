import Navigation from "./components/Navigation";
import s from "./Dashboard.module.css";

const Dashboard = () => {
  return (
    <div className={s.main}>
      <Navigation />
      <div className={s.netBalanceDiv}>
        <p>NET CALORIC BALANCE: <span>+</span>5000</p>
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
}

export default Dashboard