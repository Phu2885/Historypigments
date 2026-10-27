import "./DetailPage.css";
import { useEffect, useState, type FC } from "react";
import { useNavigate, useOutletContext, useParams } from "react-router-dom";
import { Button, Spinner } from "react-bootstrap";
import { getPigmentById, type Pigment } from "../modules/pigmentsApi";
import { ROUTES } from "../Routes";
import type { AppOutletContext } from "../AppLayout";
import { PIGMENTS_MOCK } from "../modules/mock";

const defaultImage = "/assets/images/placeholder.svg";

export const DetailPage: FC = () => {
  const [pageData, setPageData] = useState<Pigment>();
  const { id } = useParams();
  const navigate = useNavigate();
  const { cards } = useOutletContext<AppOutletContext>();

  useEffect(() => {
    if (!id) return;

    // Search in current cards list
    const fromList = cards.find((item) => String(item.id) === id);
    if (fromList) {
      setPageData(fromList);
      return;
    }

    // Call API
    getPigmentById(id)
      .then((data) => {
        if (data) {
          setPageData(data);
        } else {
          const mockItem = PIGMENTS_MOCK.find((p) => String(p.id) === id);
          setPageData(mockItem || PIGMENTS_MOCK[0]);
        }
      })
      .catch(() => {
        const mockItem = PIGMENTS_MOCK.find((p) => String(p.id) === id);
        setPageData(mockItem || PIGMENTS_MOCK[0]);
      });
  }, [id, cards]);

  const handleNext = () => {
    if (!id || !cards.length) return;
    const currentIndex = cards.findIndex((item) => String(item.id) === id);
    const nextItem = cards[currentIndex + 1];
    if (nextItem) {
      navigate(`${ROUTES.DETAIL}/${nextItem.id}`);
    }
  };

  if (!pageData) {
    return (
      <div className="detailLoader">
        <Spinner animation="border" variant="primary" />
      </div>
    );
  }

  const mediaUrl = pageData.video_url
    ? pageData.video_url.startsWith("http") || pageData.video_url.startsWith("/")
      ? pageData.video_url
      : `/assets/videos/${pageData.video_url}`
    : null;

  const imageUrl = pageData.image_url
    ? pageData.image_url.startsWith("http") || pageData.image_url.startsWith("/")
      ? pageData.image_url
      : `/assets/images/${pageData.image_url}`
    : defaultImage;

  return (
    <div className="detailPage">
      {mediaUrl ? (
        <video
          className="detailMedia"
          src={mediaUrl}
          autoPlay
          loop
          muted
          playsInline
        />
      ) : (
        <img
          className="detailMedia"
          src={imageUrl}
          alt={pageData.name}
          onError={(e) => {
            (e.target as HTMLImageElement).src = defaultImage;
          }}
        />
      )}

      <div className="detailOverlay">
        <h3 className="detailTitle">{pageData.name}</h3>
        <p className="detailDesc">
          {pageData.short_description || "Исторический пигмент искусства"}
        </p>

        <div className="opticalBadges">
          {pageData.reflection !== undefined && (
            <span className="optBadge">Отражение R: {pageData.reflection}</span>
          )}
          {pageData.absorption !== undefined && (
            <span className="optBadge">Поглощение A: {pageData.absorption}</span>
          )}
        </div>

        <Button
          className="detailNextBtn"
          variant="primary"
          onClick={handleNext}
        >
          Далее ▶
        </Button>
      </div>
    </div>
  );
};

export default DetailPage;
