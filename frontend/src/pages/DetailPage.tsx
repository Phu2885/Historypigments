import "./DetailPage.css";
import { useEffect, useState, type FC } from "react";
import { useNavigate, useOutletContext, useParams } from "react-router-dom";
import { Spinner } from "react-bootstrap";
import { getPigmentById, type Pigment } from "../modules/pigmentsApi";
import type { AppOutletContext } from "../AppLayout";
import { PIGMENTS_MOCK } from "../modules/mock";

const defaultImage = "/assets/images/placeholder.svg";

export const DetailPage: FC = () => {
  const [pageData, setPageData] = useState<Pigment>();
  const [showMore, setShowMore] = useState(false);
  const { id } = useParams();
  const navigate = useNavigate();
  const { cards } = useOutletContext<AppOutletContext>();

  useEffect(() => {
    if (!id) return;

    const fromList = cards.find((item) => String(item.id) === id);
    if (fromList) {
      setPageData(fromList);
      return;
    }

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
      navigate(`/detail/${nextItem.id}`);
    }
  };

  if (!pageData) {
    return (
      <div className="detailLoader">
        <Spinner animation="border" style={{ color: "var(--gold)" }} />
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
    <div className="feed-page">
      {mediaUrl ? (
        <video
          className="video-bg"
          src={mediaUrl}
          autoPlay
          loop
          muted
          playsInline
        />
      ) : (
        <img
          className="video-bg"
          src={imageUrl}
          alt={pageData.name}
          onError={(e) => {
            (e.target as HTMLImageElement).src = defaultImage;
          }}
        />
      )}

      <div className="feed-layout">
        <div className="feed-description-side">
          <h1 className="pigment-title">{pageData.name}</h1>

          <div className="description-container">
            <p
              className="pigment-desc"
              style={{
                display: showMore ? "block" : "-webkit-box",
                WebkitLineClamp: showMore ? "unset" : 2,
                WebkitBoxOrient: "vertical",
                overflow: showMore ? "visible" : "hidden",
              }}
            >
              {pageData.short_description || "Исторический пигмент искусства"}
            </p>
            <span
              className="desc-label"
              onClick={() => setShowMore(!showMore)}
            >
              {showMore ? "МЕНЬШЕ ▲" : "БОЛЬШЕ ▼"}
            </span>
          </div>
        </div>

        <div className="feed-params-side">
          <div className="param-group">
            <span className="param-label">ОТРАЖЕНИЕ (R)</span>
            <span className="param-value">{pageData.reflection ?? 0.60}</span>
          </div>

          <div className="param-group">
            <span className="param-label">ПОГЛОЩЕНИЕ (A)</span>
            <span className="param-value">{pageData.absorption ?? 0.40}</span>
          </div>

          <div className="action-icon">
            ❤️ <div style={{ fontSize: "12px", fontWeight: 700 }}>{pageData.likes || 0}</div>
          </div>

          <button className="next-btn-vertical" onClick={handleNext}>
            <span className="next-icon">▶</span>
            <span className="next-text">ВПРЁД</span>
          </button>
        </div>
      </div>
    </div>
  );
};

export default DetailPage;
