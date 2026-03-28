import { useNavigate } from "react-router-dom";

import dark from "@/assets/images/logo/Dark.png";
import darkAnimated from "@/assets/images/logo/DarkAnimated.gif";
import darkSmall from "@/assets/images/logo/DarkSmall.png";
// import light from "@/assets/images/logo/Light.png";
// import lightAnimated from "@/assets/images/logo/LightAnimated.gif";

export const Logo = ({
  otherStyles,
  stat,
  sm,
  dashLink,
}: {
  otherStyles?: string;
  stat?: boolean;
  sm?: boolean;
  dashLink?: boolean;
}) => {
  const navigate = useNavigate();

  const getLogo = () => {
    if (sm) {
      return darkSmall;
    }
    return stat ? dark : darkAnimated;
  };

  return (
    <img
      src={getLogo()}
      className={`cursor-pointer ${otherStyles}`}
      onClick={() => (dashLink ? navigate("/servers") : navigate("/"))}
    />
  );
};
