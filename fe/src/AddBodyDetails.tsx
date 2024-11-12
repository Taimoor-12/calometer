import { useEffect, useState, useCallback } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import Spinner from "./components/Spinner";
import { toast } from "react-toastify";
import s from "./AddBodyDetails.module.css";
import { http_post, http_get } from "./lib/http";

const AddBodyDetails = () => {
  const apiUrl = process.env.REACT_APP_API_URL;
  const navigate = useNavigate();
  const location = useLocation();

  const [loadingConfirm, setLoadingConfirm] = useState(false);
  const [loadingLogout, setLoadingLogout] = useState(false);
  const [selectedGenderOption, setSelectedGenderOption] = useState("M");
  const [selectedGoalOption, setSelectedGoalOption] = useState("L");
  const [formData, setFormData] = useState({
    age: 0,
    weight: 0.0,
    height: 0,
    gender: selectedGenderOption,
    goal: selectedGoalOption,
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;

    let parsedValue: string | number = value;
    if (name === "age" || name === "weight" || name === "height") {
      parsedValue = parseFloat(value);
    }

    if (name === "gender") {
      setSelectedGenderOption(value);
    } else if (name === "goal") {
      setSelectedGoalOption(value);
    }

    setFormData({ ...formData, [name]: parsedValue });
    console.log(name, parsedValue);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setFormData({
      ...formData,
      gender: selectedGenderOption,
      goal: selectedGoalOption,
    });
    setLoadingConfirm(true);
    console.log("submitted", formData);
    await addBodyDetailsCall();
  };

  const addBodyDetailsCall = async () => {
    try {
      const resp = await http_post(
        `${apiUrl}/api/users/body_details/add`,
        formData
      );
      const respCode = +Object.keys(resp.code)[0];
      if (respCode === 200) {
        toast.success("Successfully added body details.");
        navigate("/dashboard");
      }
    } catch (e) {
      toast.error("Something went wrong, please try again.");
    } finally {
      setLoadingConfirm(false);
    }
  };

  const handleLogout = async (e: React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    setLoadingLogout(true);
    try {
      const resp = await http_post(`${apiUrl}/api/users/logout`, {});
      const respCode = +Object.keys(resp.code)[0];
      if (respCode === 200) {
        navigate("/login", { state: { from: "addBodyDetails" } });
        toast.success("Logged out successfully");
      }
    } catch (e) {
      toast.error("Something went wrong, please try again.");
    } finally {
      setLoadingLogout(false);
    }
  };

  const doBodyDetailsExistCall = useCallback(async () => {
    const resp = await http_get(`${apiUrl}/api/users/body_details/exists`);
    console.log(resp);
    const respCode = +Object.keys(resp.code)[0];
    if (respCode === 200) {
      navigate("/dashboard");
    } else if (respCode === 401) {
      navigate("/login", { state: { from: "addBodyDetails" } });
    }
  }, [apiUrl, navigate]);

  useEffect(() => {
    if (location.state?.from === "login") {
      return;
    }

    const doBodyDetailsExist = async () => {
      await doBodyDetailsExistCall();
    };

    doBodyDetailsExist();
  }, [location.state, doBodyDetailsExistCall]);

  return (
    <div className={s.main}>
      <h1>BODY DETAILS</h1>

      <div className={s.parentDiv}>
        <div className={s.formWrapper}>
          <form onSubmit={handleSubmit}>
            <div className={s.formDiv}>
              <div className={s.inputDiv}>
                <label htmlFor="age">Age</label>
                <input
                  type="number"
                  min={20}
                  max={54}
                  id="age"
                  name="age"
                  placeholder="Enter your age"
                  onChange={handleChange}
                  required
                />
              </div>
              <div className={s.inputDiv}>
                <label htmlFor="weight">Weight</label>
                <input
                  type="number"
                  min={40}
                  max={150}
                  step={0.5}
                  id="weight"
                  name="weight"
                  placeholder="Enter weight in kg"
                  onChange={handleChange}
                  required
                />
              </div>
              <div className={s.inputDiv}>
                <label htmlFor="height">Height</label>
                <input
                  type="number"
                  min={149}
                  max={196}
                  step={0.1}
                  id="height"
                  name="height"
                  placeholder="Enter height in cm"
                  onChange={handleChange}
                  required
                />
              </div>
              <div className={s.inputDiv}>
                <label>Gender</label>
                <div className={s.radioBtnDiv}>
                  <label>
                    <div className={s.radioBtnWrapper}>
                      <input
                        type="radio"
                        value="M"
                        name="gender"
                        checked={selectedGenderOption === "M"}
                        onChange={handleChange}
                      />
                      <span>Male</span>
                    </div>
                  </label>
                  <label className={s.radioLabel}>
                    <div className={s.radioBtnWrapper}>
                      <input
                        type="radio"
                        value="F"
                        name="gender"
                        checked={selectedGenderOption === "F"}
                        onChange={handleChange}
                      />
                      <span>Female</span>
                    </div>
                  </label>
                </div>
              </div>
              <div className={s.inputDiv}>
                <label>Goal</label>
                <div className={s.radioBtnDiv}>
                  <label>
                    <div className={s.radioBtnWrapper}>
                      <input
                        type="radio"
                        value="L"
                        name="goal"
                        checked={selectedGoalOption === "L"}
                        onChange={handleChange}
                      />
                      <span>Lose</span>
                    </div>
                  </label>
                  <label className={s.radioLabel}>
                    <div className={s.radioBtnWrapper}>
                      <input
                        type="radio"
                        value="G"
                        name="goal"
                        checked={selectedGoalOption === "G"}
                        onChange={handleChange}
                      />
                      <span>Gain</span>
                    </div>
                  </label>
                </div>
              </div>
            </div>

            <div className={s.btnDiv}>
              <div className={s.loginBtn}>
                <button type="submit" disabled={loadingConfirm}>
                  {loadingConfirm ? <Spinner /> : "Confirm"}
                </button>
              </div>
              <div className={s.logoutBtn}>
                <button disabled={loadingLogout} onClick={handleLogout}>
                  {loadingLogout ? <Spinner /> : "Logout"}
                </button>
              </div>
            </div>
          </form>
        </div>
        <div className={s.imgDiv}>
          <img src="/assets/fitness_stats.svg" alt="fitness stats"></img>
        </div>
      </div>
    </div>
  );
};

export default AddBodyDetails;
