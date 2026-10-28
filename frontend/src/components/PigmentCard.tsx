import type { FC } from "react";
import { Card } from "react-bootstrap";
import "./PigmentCard.css";

const defaultImage = "/assets/images/placeholder.svg";

interface Props {
  name: string;
  short_description?: string;
  image_url?: string;
  reflection?: number;
  absorption?: number;
  likes?: number;
  onClick: () => void;
}

export const PigmentCard: FC<Props> = ({
  name,
  image_url,
  reflection,
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
    <Card className="card" onClick={onClick}>
      <div className="cardImgWrapper">
        <Card.Img
          className="card-img"
          variant="top"
          src={imgSrc}
          onError={(e) => {
            (e.target as HTMLImageElement).src = defaultImage;
          }}
        />
      </div>
      <div className="card-content">
        <Card.Title className="card-title">{name}</Card.Title>
        <div className="card-footer">
          <span>R: {reflection ?? 0.6} | A: {absorption ?? 0.4}</span>
          <span>❤️ {likes}</span>
        </div>
      </div>
    </Card>
  );
};

export default PigmentCard;
