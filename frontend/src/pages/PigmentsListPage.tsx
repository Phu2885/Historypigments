import "./PigmentsListPage.css";
import { useEffect, useState, type FC } from "react";
import { Spinner } from "react-bootstrap";
import { useNavigate, useOutletContext } from "react-router-dom";
import { getPigmentsList, type Pigment } from "../modules/pigmentsApi";
import { InputField } from "../components/InputField";
import { ROUTES } from "../Routes";
import { PigmentCard } from "../components/PigmentCard";
import type { AppOutletContext } from "../AppLayout";
import { PIGMENTS_MOCK } from "../modules/mock";

export const PigmentsListPage: FC = () => {
  const [searchValue, setSearchValue] = useState("1.0");
  const [loading, setLoading] = useState(false);
  const [pigments, setPigments] = useState<Pigment[]>(PIGMENTS_MOCK);
  const { setCards } = useOutletContext<AppOutletContext>();
  const navigate = useNavigate();

  const applyMock = () => {
    const val = parseFloat(searchValue);
    let list = PIGMENTS_MOCK;
    if (!isNaN(val)) {
      list = PIGMENTS_MOCK.filter((p) => p.absorption !== undefined && p.absorption <= val);
    }
    setPigments(list.length ? list : PIGMENTS_MOCK);
    setCards(list.length ? list : PIGMENTS_MOCK);
    setLoading(false);
  };

  const handleSearch = () => {
    setLoading(true);
    const filterVal = searchValue ? parseFloat(searchValue) : undefined;
    getPigmentsList(filterVal)
      .then((data) => {
        if (data && data.length) {
          setPigments(data);
          setCards(data);
          setLoading(false);
        } else {
          applyMock();
        }
      })
      .catch(applyMock);
  };

  useEffect(() => {
    handleSearch();
  }, []);

  const handleCardClick = (id: number) => {
    navigate(`${ROUTES.DETAIL}/${id}`);
  };

  return (
    <div className="tiles-page">
      <h2 className="section-title">Каталог пигментов</h2>

      <InputField
        value={searchValue}
        setValue={(val) => setSearchValue(val)}
        loading={loading}
        onSubmit={handleSearch}
      />

      {loading && (
        <div className="loadingBg">
          <Spinner animation="border" style={{ color: "var(--gold)" }} />
        </div>
      )}

      {!loading && (
        <div className="tiles-grid">
          {pigments.map((item) => (
            <PigmentCard
              key={item.id}
              onClick={() => handleCardClick(item.id)}
              {...item}
            />
          ))}
        </div>
      )}
    </div>
  );
};

export default PigmentsListPage;
