import React, { useContext, useEffect, useState } from "react";
import styled from "styled-components";

import Banner from "components/Banner";

import { Context } from "shared/Context";
import { Usage, UsageData } from "shared/types";
import { Link } from "react-router-dom";

type UsageKeys = keyof Usage;

const ReadableNameMap: {
  [key: string]: string;
} = {
  resource_cpu: "CPU Usage",
  resource_memory: "Memory Usage",
  clusters: "Clusters",
  users: "Users",
};

const parseToReadableString = (
  key: UsageKeys,
  current: number,
  limit: number
) => {
  switch (key) {
    case "clusters":
      return `${current} / ${limit} clusters`;
    case `resource_cpu`:
      return `${current} / ${limit} vCPU`;
    case "resource_memory":
      return `${current / 1000} / ${limit / 1000} GB`;
    case "users":
      return `${current} / ${limit} seats`;
    default:
      return `${current} / ${limit}`;
  }
};

const UpgradeChartModal: React.FC<{}> = () => {
  const { setCurrentModal, currentModalData } = useContext(Context);
  const [usage, setUsage] = useState<UsageData>(null);

  useEffect(() => {
    if (currentModalData.usage) {
      const currentUsage: UsageData = currentModalData.usage;
      setUsage(currentUsage);
    }
  }, [currentModalData?.usage]);

  if (!usage) {
    return null;
  }

  return (
    <>
      <Br />
      <Banner type="warning">
        Your project is currently exceeding its resource usage limit.
      </Banner>
      <Br />
      {usage !== null && (
        <UsageSection>
          {Object.keys(usage.current).map((key) => {
            const label = ReadableNameMap[key];
            const current = usage.current[key];
            const limit = usage.limit[key];
            const isExceeding = current > limit;
            return (
              <UsageBlock isRed={isExceeding}>
                <Label isRed={isExceeding}>{label}</Label>
                <Stat isRed={isExceeding}>
                  {parseToReadableString(key as UsageKeys, current, limit)}
                </Stat>
              </UsageBlock>
            );
          })}
        </UsageSection>
      )}
      <Helper>
        You have <b>7 days</b> to resolve this issue before your access to the
        dashboard is restricted.
      </Helper>
      <Helper>
        Have a question about billing? Email us at{" "}
        <a target="_blank" href="mailto:contact@porter.run">
          contact@porter.run
        </a>
        .
      </Helper>
      <Button
        as={Link}
        to={{
          pathname: "/project-settings",
          search: "?selected_tab=billing",
        }}
        onClick={() => setCurrentModal(null, null)}
      >
        Take me to billing
      </Button>
    </>
  );
};

export default UpgradeChartModal;

const UsageBlock = styled.div<{ isRed?: boolean }>`
  border: 1px solid ${(props) => (props.isRed ? "#ff385d" : "#ffffff55")};
  border-radius: 5px;
  padding: 18px;
`;

const Helper = styled.div`
  color: #aaaabb;
  margin-top: 12px;
`;

const Label = styled.div<{ isRed?: boolean }>`
  margin-bottom: 10px;
  font-weight: 500;
  color: ${(props) => (props.isRed ? "#ff385d" : "#ffffff55")};
`;

const Stat = styled.div<{ isRed?: boolean }>`
  font-size: 20px;
  margin-bottom: 25px;
  color: ${(props) => (props.isRed ? "#ff385d" : "#ffffff55")};
`;

const Br = styled.div`
  width: 100%;
  height: 5px;
`;

const Button = styled.button`
  height: 35px;
  background: #616feecc;
  :hover {
    background: #505edddd;
  }
  color: white;
  font-weight: 500;
  font-size: 13px;
  padding: 10px 15px;
  border-radius: 3px;
  cursor: "pointer";
  box-shadow: 0 5px 8px 0px #00000010;
  display: flex;
  flex-direction: row;
  align-items: center;
  margin-left: 10px;

  width: max-content;
  position: absolute;
  right: 20px;
  bottom: 20px;
`;

const UsageSection = styled.div`
  margin-top: 10px;
  margin-bottom: 35px;
  display: grid;
  grid-column-gap: 25px;
  grid-row-gap: 25px;
  grid-template-columns: repeat(2, minmax(200px, 1fr));
`;
