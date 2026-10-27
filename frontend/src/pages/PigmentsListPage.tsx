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
  const [searchValue, setSearchValue] = useState("");
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
    <div className="pigmentsListContainer">
      <div className="headerSection">
        <h2 className="appTitle">Alchemy of Art History</h2>
        <p className="appSubtitle">Каталог исторических пигментов</p>
        <InputField
          value={searchValue}
          setValue={(val) => setSearchValue(val)}
          loading={loading}
          onSubmit={handleSearch}
        />
      </div>

      {loading && (
        <div className="loadingBg">
          <Spinner animation="border" variant="primary" />
        </div>
      )}

      {!loading && (
        <div className="cardsList">
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
