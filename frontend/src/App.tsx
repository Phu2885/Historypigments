import { BrowserRouter, Route, Routes } from "react-router-dom";
import { DetailPage } from "./pages/DetailPage";
import PigmentsListPage from "./pages/PigmentsListPage";
import { ROUTES } from "./Routes";
import { AppLayout } from "./AppLayout";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<AppLayout />}>
          <Route path={ROUTES.LIST} element={<PigmentsListPage />} />
          <Route path={`${ROUTES.DETAIL}/:id`} element={<DetailPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
