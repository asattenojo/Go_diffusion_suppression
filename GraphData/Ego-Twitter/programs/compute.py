import networkx as nx

# エッジリストを読み込む関数
def load_twitter_ego_graph(edges_file):
    G = nx.read_edgelist(edges_file)  # NetworkXのread_edgelistを使用
    return G

def to_np_adjmat(G):
    np_adjmat = nx.to_numpy_array(G)
    # print(np_adjmat)

    return np_adjmat


# エッジリストファイルのパスを指定
edges_path = "twitter_combined/twitter_combined.txt"  # 実際のファイルパスに変更

# グラフを読み込み
G = load_twitter_ego_graph(edges_path)

#global_clustering = nx.transitivity(G)
#print(f"グローバルクラスター係数: {global_clustering}")

#average_clustering = nx.average_clustering(G)
#print(f"平均クラスター係数: {average_clustering}")

#切り抜き
print("以降　切り抜き")
import random
random.seed(1)
remaining_nodes = random.sample(list(G.nodes()), k=4000)

# 指定したノード以外を削除
G_sub = G.subgraph(remaining_nodes).copy()

G = G_sub

# 巨大連結成分を取得
largest_cc = max(nx.connected_components(G), key=len)

# 残すノードを指定
G = G.subgraph(largest_cc).copy()


# print(f"元のグラフのノード数: {G.number_of_nodes()}")
print(f"部分グラフのノード数: {G.number_of_nodes()}")


# グラフの基本情報を表示
print(G)

# グラフの可視化（matplotlibが必要）
import matplotlib.pyplot as plt
#nx.draw(G, with_labels=False, node_size=10)
#plt.show()

#print("edges:",len(G.edges))
#print("radius",nx.radius(G))
#print("diameter",nx.diameter(G))
#print("argmax:",max_index)
#print("max:",np.max(list2))
#print("ave:",np.mean(list2))
#print("min:",np.min(list2))
#print("all:",list2)

import numpy as np

# 次数分布のプロット（対数スケール）
degree_count = nx.degree_histogram(G)  # 各次数ごとのノード数
degrees = range(len(degree_count))  # 次数
frequency = np.array(degree_count) / sum(degree_count)  # 頻度

# 非ゼロのデータを抽出（対数計算が可能な部分）
nonzero_degrees = np.array([d for d in degrees if d > 0 and degree_count[d] > 0])
nonzero_frequency = np.array([frequency[d] for d in nonzero_degrees])

plt.loglog(nonzero_degrees, nonzero_frequency, marker="o", linestyle="none")
plt.title("Log-Log Degree Distribution")
plt.xlabel("Degree (log scale)")
plt.ylabel("Frequency (log scale)")
#plt.show()

from scipy.optimize import curve_fit

# フィッティング関数（べき乗則）
def power_law(x, a, b):
    return a * x ** (-b)

# フィッティングを行う
params, _ = curve_fit(power_law, nonzero_degrees, nonzero_frequency)

# フィッティング結果のプロット
plt.loglog(nonzero_degrees, nonzero_frequency, marker="o", linestyle="none", label="Data")
plt.loglog(nonzero_degrees, power_law(nonzero_degrees, *params), label=f"Fit: a={params[0]:.2f}, b={params[1]:.2f}")
plt.title("Power-Law Fit")
plt.xlabel("Degree (log scale)")
plt.ylabel("Frequency (log scale)")
plt.legend()
#plt.show()

global_clustering = nx.transitivity(G)
print(f"グローバルクラスター係数: {global_clustering}")

average_clustering = nx.average_clustering(G)
print(f"平均クラスター係数: {average_clustering}")

import pandas as pd

adj = to_np_adjmat(G)
#Golangへ無理やり持っていくように書いた
adj = adj.astype(np.int32)

print(type(adj))
# print(adj.shape)
df = pd.DataFrame(adj)

df.to_json("adj_json_egoTwitter_kirinuki.txt")
