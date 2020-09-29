from diagrams import Cluster, Diagram, Edge
from diagrams.aws.compute import EC2
from diagrams.aws.database import RDS
from diagrams.aws.network import ELB

if __name__ == "__main__":
    with Diagram("Web Service", direction="BT"):
        lb = ELB("lb")
        web = EC2("web")
        db = RDS("db")

        with Cluster("Batches"):
            batches = list(map(lambda i: EC2(f"batch_{i}"), range(1, 3)))

        lb >> web >> db << Edge(label="r/w") << batches
        db >> Edge(label="trigger", style="dotted", color="darkgreen") >> batches
