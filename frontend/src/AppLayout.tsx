import { useState, type FC } from "react";
import { Outlet } from "react-router-dom";
import type { Pigment } from "./modules/pigmentsApi";
import { TabBar } from "./components/TabBar";
import "./MobileApp.css";

export interface AppOutletContext {
  cards: Pigment[];
  setCards: (cards: Pigment[]) => void;
}

export const AppLayout: FC = () => {
  const [cards, setCards] = useState<Pigment[]>([]);

  return (
    <div className="appShell">
      <main className="appContent">
        <Outlet context={{ cards, setCards }} />
      </main>
      <TabBar cards={cards} />
    </div>
  );
};
