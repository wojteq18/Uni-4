import matplotlib.pyplot as plt
from collections import defaultdict
import numpy as np
import re

raw_data = """
n=1000, Prim: 0.4632s, Kruskal: 0.6594
n=1000, Prim: 0.3563s, Kruskal: 0.6590
n=1000, Prim: 0.2764s, Kruskal: 0.4904
n=1100, Prim: 0.3397s, Kruskal: 0.6001
n=1100, Prim: 0.3890s, Kruskal: 0.6023
n=1100, Prim: 0.3338s, Kruskal: 0.6085
n=1200, Prim: 0.4052s, Kruskal: 0.7334
n=1200, Prim: 0.3981s, Kruskal: 0.7322
n=1200, Prim: 0.3996s, Kruskal: 0.7323
n=1300, Prim: 0.4761s, Kruskal: 0.8849
n=1300, Prim: 0.4718s, Kruskal: 0.8782
n=1300, Prim: 0.5025s, Kruskal: 0.8934
n=1400, Prim: 0.5593s, Kruskal: 1.0129
n=1400, Prim: 0.5468s, Kruskal: 1.0243
n=1400, Prim: 0.5529s, Kruskal: 1.0264
n=1500, Prim: 0.6370s, Kruskal: 1.1877
n=1500, Prim: 0.6124s, Kruskal: 1.2044
n=1500, Prim: 0.6390s, Kruskal: 1.1984
n=1600, Prim: 0.6900s, Kruskal: 1.3456
n=1600, Prim: 0.7218s, Kruskal: 1.3681
n=1600, Prim: 0.7619s, Kruskal: 1.4019
n=1700, Prim: 0.8275s, Kruskal: 1.5625
n=1700, Prim: 0.8206s, Kruskal: 1.5654
n=1700, Prim: 0.8233s, Kruskal: 1.5584
n=1800, Prim: 0.9286s, Kruskal: 1.7652
n=1800, Prim: 0.9332s, Kruskal: 1.7721
n=1800, Prim: 0.9316s, Kruskal: 1.7646
n=1900, Prim: 1.0469s, Kruskal: 2.0093
n=1900, Prim: 0.9651s, Kruskal: 1.9463
n=1900, Prim: 1.0505s, Kruskal: 2.0125
n=2000, Prim: 1.0742s, Kruskal: 2.1723
n=2000, Prim: 1.0783s, Kruskal: 2.2073
n=2000, Prim: 1.1577s, Kruskal: 2.2283
n=2100, Prim: 1.2934s, Kruskal: 2.4619
n=2100, Prim: 1.2807s, Kruskal: 2.4812
n=2100, Prim: 1.2982s, Kruskal: 2.4825
n=2200, Prim: 1.4254s, Kruskal: 2.7361
n=2200, Prim: 1.4219s, Kruskal: 2.7344
n=2200, Prim: 1.4783s, Kruskal: 2.7290
n=2300, Prim: 1.6838s, Kruskal: 3.0026
n=2300, Prim: 1.5651s, Kruskal: 2.9870
n=2300, Prim: 1.5687s, Kruskal: 2.9901
n=2400, Prim: 1.7131s, Kruskal: 3.2894
n=2400, Prim: 1.7105s, Kruskal: 3.3189
n=2400, Prim: 1.7088s, Kruskal: 3.2863
n=2500, Prim: 1.8704s, Kruskal: 3.6041
n=2500, Prim: 1.9786s, Kruskal: 3.6025
n=2500, Prim: 1.8790s, Kruskal: 3.6956
n=2600, Prim: 2.0615s, Kruskal: 3.8704
n=2600, Prim: 2.0351s, Kruskal: 3.9328
n=2600, Prim: 1.8842s, Kruskal: 3.8995
n=2700, Prim: 2.0716s, Kruskal: 4.2489
n=2700, Prim: 2.2193s, Kruskal: 4.3073
n=2700, Prim: 2.7856s, Kruskal: 4.0094
n=2800, Prim: 1.4638s, Kruskal: 2.8622
n=2800, Prim: 1.4712s, Kruskal: 2.7968
n=2800, Prim: 1.4772s, Kruskal: 2.8202
n=2900, Prim: 1.5680s, Kruskal: 3.0247
n=2900, Prim: 1.5602s, Kruskal: 3.0360
n=2900, Prim: 1.5731s, Kruskal: 3.0321
n=3000, Prim: 1.6812s, Kruskal: 3.2722
n=3000, Prim: 1.7059s, Kruskal: 3.2632
n=3000, Prim: 1.7251s, Kruskal: 3.3038
n=3100, Prim: 1.7820s, Kruskal: 3.4266
n=3100, Prim: 1.7893s, Kruskal: 3.4875
n=3100, Prim: 1.8517s, Kruskal: 3.6598
n=3200, Prim: 2.0018s, Kruskal: 3.9143
n=3200, Prim: 1.9425s, Kruskal: 3.6957
n=3200, Prim: 1.9608s, Kruskal: 3.7923
n=3300, Prim: 2.0958s, Kruskal: 4.0420
n=3300, Prim: 2.1049s, Kruskal: 4.0495
n=3300, Prim: 2.1313s, Kruskal: 3.9669
n=3400, Prim: 2.2289s, Kruskal: 4.3951
n=3400, Prim: 2.1939s, Kruskal: 4.2324
n=3400, Prim: 2.1972s, Kruskal: 4.2575
n=3500, Prim: 2.3449s, Kruskal: 4.5355
n=3500, Prim: 2.3382s, Kruskal: 4.5343
n=3500, Prim: 2.3691s, Kruskal: 4.5177
n=3600, Prim: 2.5001s, Kruskal: 4.8046
n=3600, Prim: 2.4973s, Kruskal: 4.7945
n=3600, Prim: 2.4995s, Kruskal: 4.7944
n=3700, Prim: 2.6959s, Kruskal: 5.1056
n=3700, Prim: 2.7112s, Kruskal: 5.1023
n=3700, Prim: 2.7031s, Kruskal: 5.1961
n=3800, Prim: 2.8829s, Kruskal: 5.5469
n=3800, Prim: 2.8790s, Kruskal: 5.4869
n=3800, Prim: 2.8251s, Kruskal: 5.4014
n=3900, Prim: 2.9924s, Kruskal: 5.7074
n=3900, Prim: 2.9990s, Kruskal: 5.7107
n=3900, Prim: 3.0003s, Kruskal: 5.7299
n=4000, Prim: 3.1717s, Kruskal: 6.0342
n=4000, Prim: 3.2612s, Kruskal: 6.1147
n=4000, Prim: 3.3181s, Kruskal: 6.5512
n=4100, Prim: 3.5198s, Kruskal: 6.8930
n=4100, Prim: 3.4375s, Kruskal: 6.3832
n=4100, Prim: 3.4378s, Kruskal: 6.5177
n=4200, Prim: 3.6267s, Kruskal: 6.7941
n=4200, Prim: 3.6385s, Kruskal: 6.7205
n=4200, Prim: 3.6255s, Kruskal: 6.8496
n=4300, Prim: 3.9589s, Kruskal: 8.0501
n=4300, Prim: 3.9366s, Kruskal: 8.1296
n=4300, Prim: 3.8069s, Kruskal: 7.2276
n=4400, Prim: 3.9832s, Kruskal: 7.4487
n=4400, Prim: 4.0049s, Kruskal: 7.3165
n=4400, Prim: 3.9714s, Kruskal: 7.4755
n=4500, Prim: 4.1696s, Kruskal: 7.8805
n=4500, Prim: 4.1901s, Kruskal: 7.8628
n=4500, Prim: 4.5667s, Kruskal: 8.9174
n=4600, Prim: 4.5467s, Kruskal: 8.3905
n=4600, Prim: 4.4301s, Kruskal: 8.1773
n=4600, Prim: 4.4569s, Kruskal: 8.7993
n=4700, Prim: 4.8513s, Kruskal: 10.7589
n=4700, Prim: 4.7948s, Kruskal: 9.0176
n=4700, Prim: 4.6651s, Kruskal: 8.8568
n=4800, Prim: 5.1135s, Kruskal: 9.6047
n=4800, Prim: 5.0755s, Kruskal: 9.4060
n=4800, Prim: 4.8719s, Kruskal: 9.0459
n=4900, Prim: 5.1065s, Kruskal: 9.3748
n=4900, Prim: 5.0708s, Kruskal: 9.3246
n=5000, Prim: 5.2858s, Kruskal: 9.8387
n=5000, Prim: 5.3510s, Kruskal: 9.8111
n=5000, Prim: 5.3305s, Kruskal: 9.9565
n=5100, Prim: 5.5606s, Kruskal: 10.3278
n=5100, Prim: 5.6052s, Kruskal: 10.1302
n=5100, Prim: 5.6087s, Kruskal: 10.3550
n=5200, Prim: 5.8537s, Kruskal: 10.6048
n=5200, Prim: 5.8915s, Kruskal: 10.7685
n=5200, Prim: 5.8699s, Kruskal: 10.8479
n=5300, Prim: 6.1779s, Kruskal: 11.3529
n=5300, Prim: 6.2424s, Kruskal: 11.3839
n=5300, Prim: 6.1426s, Kruskal: 11.3268
n=5400, Prim: 6.4192s, Kruskal: 11.8007
n=5400, Prim: 6.4544s, Kruskal: 12.5465
n=5400, Prim: 8.5582s, Kruskal: 13.3360
n=5500, Prim: 7.0310s, Kruskal: 16.5274
n=5500, Prim: 6.6064s, Kruskal: 12.2782
n=5500, Prim: 7.0910s, Kruskal: 18.7740
n=5600, Prim: 8.4023s, Kruskal: 18.0657
n=5600, Prim: 11.8833s, Kruskal: 20.7903
n=5600, Prim: 10.3202s, Kruskal: 17.2124
n=5700, Prim: 9.7297s, Kruskal: 19.6002
n=5700, Prim: 8.7876s, Kruskal: 15.5609
n=5700, Prim: 8.4558s, Kruskal: 15.1653
n=5800, Prim: 8.1778s, Kruskal: 19.0786
n=5800, Prim: 9.0588s, Kruskal: 20.3626
n=5800, Prim: 8.3629s, Kruskal: 15.2630
n=5900, Prim: 8.8494s, Kruskal: 19.1494
n=5900, Prim: 8.4826s, Kruskal: 15.7270
n=5900, Prim: 8.4455s, Kruskal: 16.2305
n=6000, Prim: 8.8109s, Kruskal: 17.0464
n=6000, Prim: 8.9134s, Kruskal: 16.4943
n=6000, Prim: 8.8639s, Kruskal: 16.8891
n=6100, Prim: 9.3550s, Kruskal: 17.6723
n=6100, Prim: 9.3117s, Kruskal: 17.0102
n=6100, Prim: 9.6281s, Kruskal: 17.2399
n=6200, Prim: 9.5794s, Kruskal: 18.0469
n=6200, Prim: 9.7710s, Kruskal: 17.4405
n=6200, Prim: 13.3397s, Kruskal: 22.8055
n=6300, Prim: 12.9539s, Kruskal: 22.5253
n=6300, Prim: 10.7853s, Kruskal: 21.2353
n=6300, Prim: 9.8620s, Kruskal: 23.0432
n=6400, Prim: 9.9301s, Kruskal: 24.5142
"""

prim_times = defaultdict(list)
kruskal_times = defaultdict(list)

# Parsowanie danych
pattern = re.compile(r"n=(\d+), Prim: ([\d.]+)s, Kruskal: ([\d.]+)")
for line in raw_data.strip().split('\n'):
    match = pattern.match(line.strip())
    if match:
        n = int(match.group(1))
        prim = float(match.group(2))
        kruskal = float(match.group(3))
        prim_times[n].append(prim)
        kruskal_times[n].append(kruskal)

# Obliczanie średnich
n_values = sorted(prim_times.keys())
prim_avg = [np.mean(prim_times[n]) for n in n_values]
kruskal_avg = [np.mean(kruskal_times[n]) for n in n_values]

# Wykres
plt.figure(figsize=(12, 6))
plt.plot(n_values, prim_avg, label='Prim', marker='o')
plt.plot(n_values, kruskal_avg, label='Kruskal', marker='s')
plt.xlabel('Liczba wierzchołków (n)')
plt.ylabel('Średni czas wykonania (s)')
plt.title('Porównanie algorytmów Prim i Kruskal (średni czas wykonania)')
plt.legend()
plt.grid(True)
plt.tight_layout()
plt.show()