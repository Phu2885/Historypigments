import "./TabBar.css";
import type { FC } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { ROUTES } from "../Routes";
import type { Pigment } from "../modules/pigmentsApi";

interface Props {
  cards: Pigment[];
}

export const TabBar: FC<Props> = ({ cards }) => {
  const navigate = useNavigate();
  const location = useLocation();

  const isList = location.pathname === ROUTES.LIST;
  const isDetail = location.pathname.startsWith(ROUTES.DETAIL);

  const goToList = () => navigate(ROUTES.LIST);

  const goToDetail = () => {
    const firstId = cards[0]?.id;
    if (firstId) {
      navigate(`${ROUTES.DETAIL}/${firstId}`);
    } else {
      navigate(`${ROUTES.DETAIL}/3`);
    }
  };

  return (
    <nav className="bottom-nav">
      <button
        type="button"
        className={`nav-item ${isList ? "active" : ""}`}
        onClick={goToList}
      >
        <span className="nav-icon">📊</span>
        <span>Плитка</span>
      </button>

      <button
        type="button"
        className={`nav-item ${isDetail ? "active" : ""}`}
        onClick={goToDetail}
      >
        <span className="nav-icon">🎥</span>
        <span>Лента</span>
      </button>
    </nav>
  );
};

export default TabBar;
