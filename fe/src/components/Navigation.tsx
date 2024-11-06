import { useState, useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import s from "./Navigation.module.css";

const Navigation = () => {
  const navigate = useNavigate();
  const modalRef = useRef<HTMLDivElement | null>(null);
  const iconRef = useRef<HTMLDivElement | null>(null);

  const [isHamburgerIconClicked, setIsHamBurgerIconClicked] = useState(false);

  const handleLogoClick = () => {
    navigate("/dashboard");
  };

  const handleIconClick = () => {
    setIsHamBurgerIconClicked((prev) => !prev);
  };

  const handleClickOutside = (event: MouseEvent) => {
    if (
      modalRef.current &&
      iconRef.current &&
      !modalRef.current.contains(event.target as Node) &&
      !iconRef.current.contains(event.target as Node)
    ) {
      setIsHamBurgerIconClicked(false);
    }
  };

  useEffect(() => {
    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  return (
    <div className={s.navDiv}>
      <div className={s.logoDiv} onClick={handleLogoClick}>
        <img src="/assets/calometer.svg" alt="calometer" />
      </div>
      <div ref={iconRef} className={s.iconDiv} onClick={handleIconClick}>
        <img src="/assets/hamburger.svg" alt="hamburger icon" />
      </div>

      <div ref={modalRef} className={`${s.modalDiv} ${isHamburgerIconClicked ? s.modalDivFlex : s.modalDivHidden}`}>
        <div className={s.profileDiv}>
          <img src="/assets/user.svg" alt="user" />
          <span>Profile</span>
        </div>
        <div className={s.logoutDiv}>
          <img src="/assets/logout.svg" alt="logout" />
          <span>Logout</span>
        </div>
      </div>
    </div>
  );
};

export default Navigation;
