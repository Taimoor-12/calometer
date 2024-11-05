import { useEffect, useState, useCallback } from "react";
import { Link, useNavigate } from "react-router-dom";
import { http_get, http_post } from "./lib/http";
import { toast } from "react-toastify";
import Spinner from "./components/Spinner";
import s from "./Login.module.css";

interface ValidationErrors {
  username?: string;
  password?: string;
}

function Login() {
  const apiUrl = process.env.REACT_APP_API_URL;
  const navigate = useNavigate();

  const [formData, setFormData] = useState({
    username: "",
    password: "",
  });

  const [errors, setErrors] = useState<ValidationErrors>({});
  const [loading, setLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });

    const error = validateField(name, value);
    setErrors({ ...errors, [name]: error });
  };

  const validateField = (name: string, value: string): string | undefined => {
    switch (name) {
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
      default:
        break;
    }

    return undefined;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const validationErrors: ValidationErrors = {};

    if (!formData.username.trim()) {
      validationErrors.username = "Username is required";
    }

    if (!formData.password.trim()) {
      validationErrors.password = "Password is required";
    } else if (formData.password.length < 6) {
      validationErrors.password = "Password must be at least 6 characters";
    }

    setErrors(validationErrors);

    if (Object.keys(validationErrors).length === 0) {
      setLoading(true);
      console.log("Form submitted", formData);
      const body = {
        username: formData.username,
        password: formData.password,
      };

      try {
        const resp = await http_post(`${apiUrl}/api/users/login`, body);
        console.log(resp);
        const respCode = +Object.keys(resp.code)[0];
        if (respCode === 200) {
          const doBodyDetailsExist = await doBodyDetailsExistCall();
          navigate(doBodyDetailsExist ? "/dashboard" : "/addBodyDetails", {
            state: { from: "login" },
          });
          toast.success(resp.code[respCode]);
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

  const doBodyDetailsExistCall = useCallback(async (): Promise<boolean> => {
    const resp = await http_get(`${apiUrl}/api/users/body_details/exists`);
    console.log(resp);
    const respCode = +Object.keys(resp.code)[0];
    if (respCode === 200) {
      const exists: boolean = resp.data["exists"];
      return exists;
    }

    return false;
  }, [apiUrl]);

  useEffect(() => {
    const loginUser = async () => {
      try {
        const resp = await http_post(`${apiUrl}/api/users/login`, {});
        const respCode = +Object.keys(resp.code)[0];
        if (respCode === 200) {
          const doBodyDetailsExist = await doBodyDetailsExistCall();
          navigate(doBodyDetailsExist ? "/dashboard" : "/addBodyDetails", {
            state: { from: "login" },
          });
        }
      } catch (error) {
        console.error("Login failed:", error);
        toast.error("Something went wrong, please try again.");
      }
    };

    loginUser();
  }, [apiUrl, navigate, doBodyDetailsExistCall]);

  return (
    <div className={s.main}>
      <h1>SIGNIN</h1>

      <div className={s.parentDiv}>
        <div className={s.formDiv}>
          <form onSubmit={handleSubmit}>
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

            <div className={s.loginBtn}>
              <button type="submit" disabled={loading}>
                {loading ? <Spinner /> : "Login"}
              </button>
            </div>
          </form>

          <p className={s.signupPrompt}>
            Don't have an account?{" "}
            <span>
              <Link to="/signup">Sign up here</Link>
            </span>
          </p>
        </div>
        <div className={s.imgDiv}>
          <img src="/assets/login_illustration.svg" alt="login illustration" />
        </div>
      </div>
    </div>
  );
}

export default Login;
