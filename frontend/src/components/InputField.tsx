import type { FC } from "react";
import { Button } from "react-bootstrap";
import "./InputField.css";

interface Props {
  value: string;
  setValue: (value: string) => void;
  onSubmit: () => void;
  loading?: boolean;
  placeholder?: string;
  buttonTitle?: string;
}

export const InputField: FC<Props> = ({
  value,
  setValue,
  onSubmit,
  loading,
  placeholder = "Поглощение (A <= ...)",
  buttonTitle = "Искать",
}) => (
  <div className="inputField">
    <input
      type="text"
      value={value}
      placeholder={placeholder}
      onChange={(e) => setValue(e.target.value)}
      onKeyDown={(e) => e.key === "Enter" && onSubmit()}
    />
    <Button disabled={loading} variant="primary" onClick={onSubmit}>
      {buttonTitle}
    </Button>
  </div>
);

export default InputField;
