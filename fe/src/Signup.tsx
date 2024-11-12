import { useState, useEffect } from "react";
import { Link, useNavigate } from "react-router-dom";
import { http_post } from "./lib/http";
import Spinner from "./components/Spinner";
import { toast } from "react-toastify";
import s from "./Signup.module.css";

interface ValidationErrors {
  fullName?: string;
  username?: string;
  password?: string;
  confirmPassword?: string;
}

function Signup() {
  const apiUrl = process.env.REACT_APP_API_URL;
  const navigate = useNavigate();

  const [formData, setFormData] = useState({
    fullName: "",
    username: "",
    password: "",
    confirmPassword: "",
  });

  const [errors, setErrors] = useState<ValidationErrors>({});
  const [loading, setLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });

    const error = validateField(name, value);
    setErrors({ ...errors, [name]: error });
  };

  const validateField = (name: string, value: string): string | undefined => {
    switch (name) {
      case "fullName":
        if (!value.trim()) {
          return "Full name is required";
        } else if (!/^[A-Za-z\s]+$/.test(value)) {
          return "Full name must contain only letters";
        }
        break;
      case "username":
        if (!value.trim()) {
          return "Username is required";
        }
        break;
      case "password":
        if (!value.trim()) {
          return "Password is required";
        } else if (value.length < 6) {
          return "Password must be at least 6 characters";
        }
        break;
      case "confirmPassword":
        if (!value.trim()) {
          return "Confirm password is required";
        } else if (value !== formData.password) {
          return "Passwords do not match";
        }
        break;
      default:
        break;
    }

    return undefined;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const validationErrors: ValidationErrors = {};

    if (!formData.fullName.trim()) {
      validationErrors.fullName = "Full name is required";
    } else if (!/^[A-Za-z\s]+$/.test(formData.fullName)) {
      validationErrors.fullName = "Full name must contain only letters";
    }

    if (!formData.username.trim()) {
      validationErrors.username = "Username is required";
    }

    if (!formData.password.trim()) {
      validationErrors.password = "Password is required";
    } else if (formData.password.length < 6) {
      validationErrors.password = "Password must be at least 6 characters";
    }

    if (!formData.confirmPassword.trim()) {
      validationErrors.confirmPassword = "Confirm password is required";
    } else if (formData.password !== formData.confirmPassword) {
      validationErrors.confirmPassword = "Passwords do not match";
    }

    setErrors(validationErrors);

    if (Object.keys(validationErrors).length === 0) {
      setLoading(true);
      console.log("Form submitted", formData);
      const body = {
        name: formData.fullName,
        username: formData.username,
        password: formData.password,
      };

      try {
        const resp = await http_post(`${apiUrl}/api/users/signup`, body);
        console.log(resp);
        const respCode = +Object.keys(resp.code)[0];
        if (respCode === 200) {
          navigate("/login");
          toast.success("Signed up successfully. Please log in.");
        } else {
          toast.error(resp.code[respCode]);
        }
      } catch (error) {
        console.error("Login failed:", error);
        toast.error("Something went wrong, please try again.");
      } finally {
        setLoading(false);
      }
    }
  };

  useEffect(() => {
    const signupUser = async () => {
      try {
        const resp = await http_post(`${apiUrl}/api/users/signup`, {});
        console.log(resp);
        const respCode = +Object.keys(resp.code)[0];
        if (respCode === 200) {
          navigate("/login");
        }
      } catch (error) {
        console.error("Signup failed:", error);
        toast.error("Something went wrong, please try again.");
      }
    };

    signupUser();
  }, [apiUrl, navigate]);

  return (
    <div className={s.main}>
      <h1>SIGNUP</h1>

      <div className={s.parentDiv}>
        <div className={s.formDiv}>
          <form onSubmit={handleSubmit}>
            <div className={s.inputDiv}>
              <p>Full Name</p>
              <input
                type="text"
                name="fullName"
                value={formData.fullName}
                placeholder="Enter your fullname"
                onChange={handleChange}
              />
              {errors.fullName && <p className={s.error}>{errors.fullName}</p>}
            </div>

            <div className={s.inputDiv}>
              <p>Username</p>
              <input
                type="text"
                name="username"
                value={formData.username}
                placeholder="Enter your username"
                onChange={handleChange}
              />
              {errors.username && <p className={s.error}>{errors.username}</p>}
            </div>

            <div className={s.inputDiv}>
              <p>Password</p>
              <div className={s.passwordContainer}>
                <input
                  type={showPassword ? "text" : "password"}
                  name="password"
                  value={formData.password}
                  placeholder="Enter your password"
                  onChange={handleChange}
                />
                <div
                  className={s.showPasswordDiv}
                  onClick={() => setShowPassword(!showPassword)}
                >
                  {showPassword ? (
                    <img src="/assets/eye.svg" alt="eye" />
                  ) : (
                    <img src="/assets/eye-off.svg" alt="eye-off" />
                  )}
                </div>
              </div>
              {errors.password && <p className={s.error}>{errors.password}</p>}
            </div>

            <div className={s.inputDiv}>
              <p>Confirm Password</p>
              <div className={s.passwordContainer}>
                <input
                  type={showConfirmPassword ? "text" : "password"}
                  name="confirmPassword"
                  value={formData.confirmPassword}
                  placeholder="Enter the password again"
                  onChange={handleChange}
                />
                <div
                  className={s.showPasswordDiv}
                  onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                >
                  {showConfirmPassword ? (
                    <img src="/assets/eye.svg" alt="eye" />
                  ) : (
                    <img src="/assets/eye-off.svg" alt="eye-off" />
                  )}
                </div>
              </div>
              {errors.confirmPassword && (
                <p className={s.error}>{errors.confirmPassword}</p>
              )}
            </div>

            <div className={s.registerBtn}>
              <button type="submit" disabled={loading}>
                {loading ? <Spinner /> : "Register"}
              </button>
            </div>
          </form>

          <p className={s.loginPrompt}>
            Already have an account?{" "}
            <span>
              <Link to="/login">Login here</Link>
            </span>
          </p>
        </div>
        <div className={s.imgDiv}>
          <img src="/assets/healthy_habit.svg" alt="healthy habit" />
        </div>
      </div>
    </div>
  );
}

export default Signup;
