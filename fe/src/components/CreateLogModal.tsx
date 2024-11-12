import { FC, useState, useRef } from "react";
import { http_post } from "../lib/http";
import { toast } from "react-toastify";
import Spinner from "./Spinner";
import s from "./CreateLogModal.module.css";

interface CreateLogModalProps {
  toggleModalOpen: () => void;
}

const CreateLogModal: FC<CreateLogModalProps> = ({ toggleModalOpen }) => {
  const apiUrl = process.env.REACT_APP_API_URL;

  const [selectedDate, setSelectedDate] = useState<string>(
    new Date().toISOString().split("T")[0] // Initialize with today's date
  );
  const [isLoading, setIsLoading] = useState(false);


  const dateInputRef = useRef<HTMLInputElement>(null);

  const handleInputClick = () => {
    dateInputRef.current?.showPicker?.();
    dateInputRef.current?.focus();
  };

  const handleDateChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setSelectedDate(event.target.value);
  };

  const handleCreateClick = async () => {
    setIsLoading(true);

    // Convert selected date to ISO format for Golang backend
    const formattedDate = new Date(selectedDate).toISOString(); // Output in `YYYY-MM-DDTHH:MM:SSZ` format
    console.log("Formatted Date for API:", formattedDate);

    const body = {
      log_date: formattedDate
    }

    try {
      const resp = await http_post(`${apiUrl}/api/users/log/create`, body)
      const respCode = +Object.keys(resp.code)[0]
      if (respCode === 400 || respCode === 409) {
        toast.error(resp.code[respCode])
      } else if (respCode === 200) {
        toast.success("Log has been created.")
      }
    } catch (e) {
      toast.error("Something went wrong, please try again.");
    } finally {
      setIsLoading(false);
    }
    
  };

  return (
    <div className={s.container}>
      <div onClick={toggleModalOpen} className={s.overlay}></div>
      <div className={s.contentWrapper}>
        <img onClick={toggleModalOpen} src="/assets/cross.svg" alt="close" />
        <div className={s.contentDiv}>
          <h1>CREATE LOG</h1>
          <div className={s.inputDiv}>
            <p>Date</p>
            <input
              type="date"
              name="date"
              min="2024-01-01"
              max="2025-12-31"
              defaultValue={selectedDate}
              ref={dateInputRef}
              onClick={handleInputClick}
              onChange={handleDateChange}
              className={s.dateInput}
            />
          </div>
          <div className={s.btnDiv} onClick={handleCreateClick}>
            <button>
              {isLoading ? <Spinner /> : "Create"}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CreateLogModal;
