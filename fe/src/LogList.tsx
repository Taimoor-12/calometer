import { FC } from "react";
import s from "./LogList.module.css";

interface LogListProps {
  setIsTileClicked: React.Dispatch<React.SetStateAction<boolean>>;
}

const LogList: FC<LogListProps> = ({ setIsTileClicked }) => {
  const handleArrowBackClick = () => {
    setIsTileClicked(false);
  }
  
  return (
    <div className={s.main}>
      <img
        src="/assets/arrow_back.svg"
        alt="arrow back"
        className={s.arrowBack}
        onClick={handleArrowBackClick}
      />
    </div>
  );
};

export default LogList;
