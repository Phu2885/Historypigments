import type { FC } from "react";
import { Card } from "react-bootstrap";
import "./PigmentCard.css";

const defaultImage = "/assets/images/placeholder.svg";

interface Props {
  name: string;
  short_description?: string;
  image_url?: string;
  absorption?: number;
  likes?: number;
  onClick: () => void;
}

export const PigmentCard: FC<Props> = ({
  name,
  short_description,
  image_url,
  absorption,
  likes = 0,
  onClick,
}) => {
  const imgSrc = image_url
    ? image_url.startsWith("http") || image_url.startsWith("/")
      ? image_url
      : `/assets/images/${image_url}`
    : defaultImage;

  return (
    <Card className="pigmentCard" onClick={onClick}>
      <div className="pigmentImgWrapper">
        <Card.Img
          className="pigmentCardImg"
          variant="top"
          src={imgSrc}
          onError={(e) => {
            (e.target as HTMLImageElement).src = defaultImage;
          }}
        />
        {absorption !== undefined && (
          <span className="absorptionBadge">A ≤ {absorption}</span>
        )}
      </div>
      <Card.Body className="pigmentCardBody">
        <Card.Title className="pigmentCardTitle">{name}</Card.Title>

        <Card.Text className="pigmentCardDesc">
          {short_description || "Исторический пигмент искусства"}
        </Card.Text>

        <div className="pigmentCardFooter">
          <span className="pigmentLikes">❤️ {likes}</span>
        </div>
      </Card.Body>
    </Card>
  );
};

export default PigmentCard;
