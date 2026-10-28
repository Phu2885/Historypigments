import type { FC } from "react";
import { Button } from "react-bootstrap";
import "./InputField.css";

interface Props {
  value: string;
  setValue: (value: string) => void;
  onSubmit: () => void;
  loading?: boolean;
}

export const InputField: FC<Props> = ({
  value,
  setValue,
  onSubmit,
  loading,
}) => {
  const currentVal = value ? parseFloat(value) : 1.0;

  return (
    <div className="inputField">
      <div className="inputFieldHeader">
        <span className="inputFieldLabel">ФИЛЬТР ПОГЛОЩЕНИЯ (A ≤ {currentVal.toFixed(2)})</span>
      </div>

      <input
        type="range"
        min="0.1"
        max="1.0"
        step="0.02"
        className="filterSlider"
        value={currentVal}
        onChange={(e) => setValue(e.target.value)}
      />

      <Button
        disabled={loading}
        className="searchBtn"
        onClick={onSubmit}
      >
        {loading ? "ПОИСК..." : "ПРИМЕНИТЬ ФИЛЬТР"}
      </Button>
    </div>
  );
};

export default InputField;
